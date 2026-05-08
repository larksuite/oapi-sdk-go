package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/channel"
	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

func main() {
	// 支持命令行参数
	var appID string
	var appSecret string
	var dmMode string
	var respondToMentionAll bool

	flag.StringVar(&appID, "app_id", "", "Lark App ID (也可以通过 APP_ID 环境变量传入)")
	flag.StringVar(&appSecret, "app_secret", "", "Lark App Secret (也可以通过 APP_SECRET 环境变量传入)")
	flag.StringVar(&dmMode, "dm_mode", "open", "单聊策略模式: open(默认), disabled, allowlist")
	flag.BoolVar(&respondToMentionAll, "respond_all", true, "群聊策略: 是否响应 @所有人 (默认 true)")
	flag.Parse()

	// 从环境变量获取配置作 fallback
	if appID == "" {
		appID = os.Getenv("APP_ID")
	}
	if appSecret == "" {
		appSecret = os.Getenv("APP_SECRET")
	}

	if appID == "" || appSecret == "" {
		fmt.Println("启动失败: 缺少必要的参数")
		fmt.Println("使用方式: go run main.go -app_id=cli_xxx -app_secret=xxx")
		fmt.Println("或者设置环境变量: APP_ID=cli_xxx APP_SECRET=xxx go run main.go")
		return
	}

	// 1. 创建 Lark API Client
	client := lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithSource("oapi-sdk-go-sample"), // 追加 User-Agent 追踪
	)

	// 2. 创建 EventDispatcher，用于 WebSocket 的事件路由
	eventHandler := dispatcher.NewEventDispatcher("", "")

	// 3. 创建 WebSocket Client，并将 EventHandler 注入
	wsClient := larkws.NewClient(appID, appSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
		larkws.WithOnReady(func() {
			fmt.Println("[WS] 连接已就绪!")
		}),
		larkws.WithOnError(func(err error) {
			fmt.Printf("[WS] 连接发生错误: %v\n", err)
		}),
		larkws.WithOnReconnecting(func() {
			fmt.Println("[WS] 正在尝试重连...")
		}),
		larkws.WithOnReconnected(func() {
			fmt.Println("[WS] 重连成功!")
		}),
	)

	// 4. 创建 Channel 抽象实例，将 client 和 wsClient 传入
	// 默认由外部启动参数控制安全策略：
	// -respond_all=true (测试 TC-311 @所有人)
	// -dm_mode=disabled (测试 TC-502 单聊拒绝)
	ch := channel.NewChannel(client, wsClient, types.WithPolicyConfig(types.PolicyConfig{
		RespondToMentionAll: &respondToMentionAll,
		DMMode:              dmMode,
	}))

	// 5. 注册消息处理逻辑
	ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
		fmt.Printf("收到来自 %s 的消息: %s\n", msg.UserID, msg.Content)

		// 将转换后的归一化消息 (NormalizedMessage) 输出为 JSON 字符串
		if normJSON, err := json.MarshalIndent(msg, "", "  "); err == nil {
			fmt.Printf("📝 NormalizedMessage: %s\n", string(normJSON))
		}

		// 打印原始事件，便于调试排查和比对结构
		if rawJSON, err := json.Marshal(msg.RawEvent); err == nil {
			fmt.Printf("📦 RawEvent: %s\n", string(rawJSON))
		}

		msgType := msg.RawContentType

		// TC-301
		if msgType == "text" {
			if msg.Content != "" && msg.RawContentType == "text" {
				fmt.Printf("✅ [TC-301 Passed] Received Text Message (Content: %s)\n", msg.Content)
			} else {
				fmt.Printf("❌ [TC-301 Failed] Content is empty or RawContentType mismatch\n")
			}
		}
		// TC-302
		if msgType == "image" {
			if (strings.Contains(msg.Content, "![image]") || strings.Contains(msg.Content, "![image](")) && len(msg.Resources) > 0 && msg.Resources[0].Type == "image" && msg.RawContentType == "image" {
				fmt.Printf("✅ [TC-302 Passed] Received Image Message (ImageKey: %s)\n", msg.Resources[0].FileKey)
			} else {
				fmt.Printf("❌ [TC-302 Failed] Missing image tag, resource, or RawContentType mismatch\n")
			}
		}
		// TC-303
		if msgType == "post" {
			if msg.Content != "" && msg.RawContentType == "post" {
				fmt.Printf("✅ [TC-303 Passed] Received Post Message (Converted Markdown: %s)\n", msg.Content)
			} else {
				fmt.Printf("❌ [TC-303 Failed] Post content is empty or RawContentType mismatch\n")
			}
		}
		// TC-304
		if msgType == "file" {
			if strings.Contains(msg.Content, "<file") && len(msg.Resources) > 0 && msg.Resources[0].Type == "file" && msg.RawContentType == "file" {
				fmt.Printf("✅ [TC-304 Passed] Received File Message (FileKey: %s)\n", msg.Resources[0].FileKey)
			} else {
				fmt.Printf("❌ [TC-304 Failed] Missing file tag, resource, or RawContentType mismatch\n")
			}
		}
		// TC-305
		if msgType == "audio" {
			if strings.Contains(msg.Content, "<audio") && len(msg.Resources) > 0 && msg.Resources[0].Type == "audio" && msg.RawContentType == "audio" {
				fmt.Printf("✅ [TC-305 Passed] Received Audio Message (AudioKey: %s)\n", msg.Resources[0].FileKey)
			} else {
				fmt.Printf("❌ [TC-305 Failed] Missing audio tag, resource, or RawContentType mismatch\n")
			}
		}
		// TC-306
		if msgType == "media" {
			if strings.Contains(msg.Content, "<video") && len(msg.Resources) > 0 && msg.Resources[0].Type == "video" && msg.RawContentType == "media" {
				fmt.Printf("✅ [TC-306 Passed] Received Video Message (VideoKey: %s)\n", msg.Resources[0].FileKey)
			} else {
				fmt.Printf("❌ [TC-306 Failed] Missing video tag, resource, or RawContentType mismatch\n")
			}
		}
		// TC-307
		if msgType == "share_chat" || msgType == "share_user" {
			if strings.Contains(msg.Content, "<share_") && (msg.RawContentType == "share_chat" || msg.RawContentType == "share_user") {
				fmt.Printf("✅ [TC-307 Passed] Received Share Card (Type: %s)\n", msgType)
			} else {
				fmt.Printf("❌ [TC-307 Failed] Missing share tag in content or RawContentType mismatch\n")
			}
		}
		// TC-308
		if msgType == "interactive" {
			if msg.Content != "" && msg.RawContentType == "interactive" {
				fmt.Printf("✅ [TC-308 Passed] Received Interactive Card\n")
			} else {
				fmt.Printf("❌ [TC-308 Failed] Card content extraction failed or RawContentType mismatch\n")
			}
		}
		// TC-309
		if msgType == "merge_forward" {
			if msg.Content != "" && msg.RawContentType == "merge_forward" {
				fmt.Printf("✅ [TC-309 Passed] Received Merge Forward Message\n")
			} else {
				fmt.Printf("❌ [TC-309 Failed] Missing merge forward content or RawContentType mismatch\n")
			}
		}

		// TC-310
		if msg.MentionedBot {
			if len(msg.Mentions) > 0 {
				fmt.Printf("✅ [TC-310 Passed] Bot was mentioned (Mentions count: %d)\n", len(msg.Mentions))
			} else {
				fmt.Printf("❌ [TC-310 Failed] MentionedBot is true but mentions array is empty\n")
			}
		}
		// TC-311
		if msg.MentionAll {
			fmt.Printf("✅ [TC-311 Passed] Received @all message\n")
		}
		// TC-312
		if msg.ChatType == "p2p" || msg.ChatType == "group" {
			fmt.Printf("✅ [TC-312 Passed] Received message from %s chat\n", msg.ChatType)
		} else if msg.ChatType != "" {
			fmt.Printf("❌ [TC-312 Failed] Invalid ChatType %s\n", msg.ChatType)
		}

		// 尽可能贴近真实的指令解析方式：处理群聊中包含 @机器人的 前缀
		// 飞书群聊里 @机器人 时，纯文本内容通常会带上 "@_user_1" 或者真实的机器名字
		cmd := strings.TrimSpace(msg.Content)

		// 简单粗暴的正则/前缀清理：如果有以 @ 开头的文本（通常就是@机器人），直接把它剥离掉
		if strings.HasPrefix(cmd, "@") {
			parts := strings.SplitN(cmd, " ", 2)
			if len(parts) == 2 {
				cmd = strings.TrimSpace(parts[1])
			} else {
				cmd = "" // 只有 @机器人，没有带指令
			}
		}

		switch {
		case strings.HasPrefix(cmd, "/policy"):
			// 动态修改安全策略，方便手工测试拦截 (TC-501~TC-505)
			// 例如: /policy dm_mode=disabled, respond_all=false, require_mention=true
			parts := strings.SplitN(cmd, " ", 2)
			if len(parts) == 2 {
				args := strings.Split(parts[1], ",")
				pol := ch.GetPolicy()
				bTrue := true
				bFalse := false
				for _, arg := range args {
					arg = strings.TrimSpace(arg)
					kv := strings.SplitN(arg, "=", 2)
					if len(kv) == 2 {
						k := strings.TrimSpace(kv[0])
						v := strings.TrimSpace(kv[1])
						switch k {
						case "dm_mode":
							pol.DMMode = v
						case "respond_all":
							if v == "true" {
								pol.RespondToMentionAll = &bTrue
							} else {
								pol.RespondToMentionAll = &bFalse
							}
						case "require_mention":
							if v == "true" {
								pol.RequireMention = &bTrue
							} else {
								pol.RequireMention = &bFalse
							}
						case "group_allowlist":
							if v == "empty" || v == "" {
								pol.GroupAllowlist = []string{}
							} else {
								pol.GroupAllowlist = strings.Split(v, "|")
							}
						case "dm_allowlist":
							if v == "empty" || v == "" {
								pol.DMAllowlist = []string{}
							} else {
								pol.DMAllowlist = strings.Split(v, "|")
							}
						}
					}
				}
				ch.UpdatePolicy(pol)
				// 打印更新后的配置
				newPolJSON, _ := json.MarshalIndent(ch.GetPolicy(), "", "  ")
				_, _ = ch.Send(ctx, &types.SendInput{
					ReceiveID:      msg.ChatID,
					ReplyMessageID: msg.MessageID,
					Text:           "安全策略已更新: \n" + string(newPolJSON),
				})
			} else {
				// 打印当前配置
				curPolJSON, _ := json.MarshalIndent(ch.GetPolicy(), "", "  ")
				_, _ = ch.Send(ctx, &types.SendInput{
					ReceiveID:      msg.ChatID,
					ReplyMessageID: msg.MessageID,
					Text:           "当前安全策略: \n" + string(curPolJSON) + "\n\n更新示例:\n/policy dm_mode=disabled\n/policy respond_all=false\n/policy require_mention=true\n/policy group_allowlist=oc_xxx|oc_yyy",
				})
			}
		case strings.HasPrefix(cmd, "/stream"):
			// 测试 Stream (流式响应)
			stream, err := ch.Stream(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Markdown:  "正在思考中...\n",
			})
			if err != nil {
				return err
			}

			// 模拟慢速生成文本
			for i := 1; i <= 5; i++ {
				stream.Append(ctx, fmt.Sprintf("这是第 %d 行流式文本...\n", i))
				time.Sleep(1 * time.Second)
			}

			stream.Append(ctx, "\n**完成！**")
			stream.Close(ctx)

		case strings.HasPrefix(cmd, "/markdown"):
			// 测试富文本发送
			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Title:     "测试富文本",
				Markdown:  "这是第一段\n\n这是第二段，带有 [链接](https://larksuite.com)\n\n<at user_id=\"all\">所有人</at>",
			})
			if err != nil {
				return err
			}

		case strings.HasPrefix(cmd, "/cardstream"):
			// 测试 CardStream (流式卡片更新)
			stream, err := ch.Stream(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Card: `{
					"elements": [
						{
							"tag": "div",
							"text": {
								"content": "⌛️ 任务正在处理中... 0%",
								"tag": "lark_md"
							}
						}
					]
				}`,
			})
			if err != nil {
				return err
			}

			// 模拟多步排队更新卡片状态
			for i := 1; i <= 5; i++ {
				progress := i * 20
				updatedCard := fmt.Sprintf(`{
					"elements": [
						{
							"tag": "div",
							"text": {
								"content": "🚀 任务正在处理中... %d%%",
								"tag": "lark_md"
							}
						}
					]
				}`, progress)
				stream.UpdateCard(ctx, updatedCard)
				time.Sleep(1 * time.Second)
			}

			// 最后更新为完成状态
			stream.UpdateCard(ctx, `{
				"elements": [
					{
						"tag": "div",
						"text": {
							"content": "✅ 任务已处理完成！100%",
							"tag": "lark_md"
						}
					}
				]
			}`)
			stream.Close(ctx)

		case strings.HasPrefix(cmd, "/card"):
			// 测试卡片发送
			cardJSON := `{
				"elements": [
					{
						"tag": "div",
						"text": {
							"content": "这是一个交互卡片示例",
							"tag": "lark_md"
						}
					},
					{
						"tag": "action",
						"actions": [
							{
								"tag": "button",
								"text": {
									"content": "点我测试",
									"tag": "plain_text"
								},
								"type": "primary",
								"value": {
									"action": "test_button_click"
								}
							}
						]
					}
				]
			}`
			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Card:      cardJSON,
			})
			if err != nil {
				return err
			}

		case strings.HasPrefix(cmd, "/mention"):
			// 测试 @ 提及组合器
			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Text:      "这是通过 Mentions 字段组合的提及消息",
				Mentions: []types.Mention{
					{UserID: msg.UserID}, // 自动 @ 发送者
				},
			})
			if err != nil {
				return err
			}

		case strings.HasPrefix(cmd, "/file"):
			// 测试文件发送 (内置 Uploader 自动上传并发送)
			tmpFile := filepath.Join(os.TempDir(), "oapi_test_file.txt")
			_ = os.WriteFile(tmpFile, []byte("Hello, this is a test file auto-uploaded by oapi-sdk-go channel!"), 0644)
			defer os.Remove(tmpFile) // 测试完清理

			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				FilePath:  tmpFile,
			})
			if err != nil {
				return err
			}

		case strings.HasPrefix(cmd, "/image"):
			// 测试图片发送 (内置 Uploader 自动上传并发送)
			tmpImg := filepath.Join(os.TempDir(), "oapi_test_image.png")
			// 生成一个极小的 1x1 透明 PNG
			pngData, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=")
			_ = os.WriteFile(tmpImg, pngData, 0644)
			defer os.Remove(tmpImg) // 测试完清理

			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				ImagePath: tmpImg,
			})
			if err != nil {
				return err
			}

		case strings.HasPrefix(cmd, "/post"):
			// 测试发送富文本 post
			postJSON := `{"zh_cn": {"title": "我是一个标题", "content": [[{"tag": "text", "text": "我是富文本内容"}]]}}`
			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Post:      postJSON,
			})
			if err != nil {
				return err
			}

		case strings.HasPrefix(cmd, "/sharechat"):
			// 测试发送群名片
			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID:   msg.ChatID,
				ShareChatID: msg.ChatID,
			})
			if err != nil {
				return err
			}

		case strings.HasPrefix(cmd, "/shareuser"):
			// 测试发送个人名片，支持 /shareuser xxx@example.com 或 /shareuser -email=xxx@example.com
			parts := strings.SplitN(cmd, " ", 2)
			targetUserID := msg.UserID // 默认分享发送者自己

			if len(parts) == 2 {
				email := strings.TrimSpace(parts[1])
				// 如果飞书客户端自动把邮箱转成了超链接，格式可能是 [xxx@example.com](mailto:xxx@example.com)
				// 我们需要用简单的正则或字符串操作把它提取出来
				if strings.Contains(email, "](mailto:") {
					start := strings.Index(email, "mailto:") + 7
					end := strings.Index(email[start:], ")")
					if start > 6 && end > 0 {
						email = email[start : start+end]
					}
				}
				if strings.HasPrefix(email, "-email=") {
					email = strings.TrimPrefix(email, "-email=")
				}

				if email != "" {
					req := larkcontact.NewBatchGetIdUserReqBuilder().
						UserIdType("open_id").
						Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
							Emails([]string{email}).
							Build()).
						Build()
					resp, err := client.Contact.V3.User.BatchGetId(ctx, req)
					if err == nil && resp != nil && resp.Success() && resp.Data != nil && len(resp.Data.UserList) > 0 && resp.Data.UserList[0] != nil && resp.Data.UserList[0].UserId != nil {
						targetUserID = *resp.Data.UserList[0].UserId
					} else {
						// 查询失败，给个提示并退回
						_, _ = ch.Send(ctx, &types.SendInput{
							ReceiveID: msg.ChatID,
							Text:      fmt.Sprintf("未找到该邮箱对应的用户或查询失败: %s", email),
						})
						return nil
					}
				}
			}

			// 如果 targetUserID 不是以 ou_ 开头的标准 open_id，通常在自建应用测试桩或内部环境中会报错 230001 invalid user_id
			// 因为飞书 API 要求 share_user 的值必须是合法的 open_id。这里给一个兜底提示：
			if !strings.HasPrefix(targetUserID, "ou_") {
				_, _ = ch.Send(ctx, &types.SendInput{
					ReceiveID: msg.ChatID,
					Text:      fmt.Sprintf("⚠️ 无法分享名片：当前提取的 UserID (%s) 不是合法的 open_id，请使用邮箱参数调用，例如：/shareuser xxx@example.com", targetUserID),
				})
				return nil
			}

			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID:   msg.ChatID,
				ShareUserID: targetUserID,
			})
			if err != nil {
				return err
			}

		default:
			// 测试普通 Send (回显消息)
			normJSONStr := ""
			if normJSON, err := json.MarshalIndent(msg, "", "  "); err == nil {
				normJSONStr = string(normJSON)
			}

			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID:      msg.ChatID,
				ReplyMessageID: msg.MessageID,
				Text:           "Echo: \n" + normJSONStr + "\n\n💡 支持指令:\n/policy\n/stream\n/markdown\n/card\n/cardstream\n/mention\n/file\n/image\n/post\n/sharechat\n/shareuser",
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	// 6. 注册边缘事件与生命周期
	ch.OnBotAdded(func(ctx context.Context, event *types.BotAddedEvent) error {
		if normJSON, err := json.MarshalIndent(event, "", "  "); err == nil {
			fmt.Printf("🤖 BotAddedEvent: %s\n", string(normJSON))
		}
		if event.ChatID != "" && event.ChatName != "" && event.UserID != "" {
			fmt.Printf("✅ [TC-316 Passed] 机器人被添加到群聊: %s (ChatID: %s) by UserID: %s\n", event.ChatName, event.ChatID, event.UserID)
		} else {
			fmt.Printf("❌ [TC-316 Failed] Bot added event missing chat or operator info\n")
		}
		return nil
	})

	ch.OnReaction(func(ctx context.Context, event *types.ReactionEvent) error {
		if normJSON, err := json.MarshalIndent(event, "", "  "); err == nil {
			fmt.Printf("👍 ReactionEvent: %s\n", string(normJSON))
		}
		if event.ReactionType == "" || event.UserID == "" {
			fmt.Printf("❌ [TC-314/315 Failed] Reaction missing type or user ID\n")
			return nil
		}
		if event.Action == "added" {
			fmt.Printf("✅ [TC-314 Passed] 收到表情回应: %s, User: %s, messageID: %s\n", event.ReactionType, event.UserID, event.MessageID)
		} else if event.Action == "removed" {
			fmt.Printf("✅ [TC-315 Passed] 表情回应被移除: %s, User: %s, messageID: %s\n", event.ReactionType, event.UserID, event.MessageID)
		} else {
			fmt.Printf("❌ [TC-314/315 Failed] Unknown reaction action %s\n", event.Action)
		}
		return nil
	})

	ch.OnComment(func(ctx context.Context, event *types.CommentEvent) error {
		if normJSON, err := json.MarshalIndent(event, "", "  "); err == nil {
			fmt.Printf("💬 CommentEvent: %s\n", string(normJSON))
		}
		if event.CommentID != "" && event.FileToken != "" && event.Operator.OpenID != "" {
			fmt.Printf("✅ [TC-317 Passed] 收到文档评论/回复: CommentID: %s, FileToken: %s, Operator: %s, MentionedBot: %v\n", event.CommentID, event.FileToken, event.Operator.OpenID, event.MentionedBot)
		} else {
			fmt.Printf("❌ [TC-317 Failed] Comment event missing commentId, fileToken or operator\n")
		}
		return nil
	})

	ch.OnReject(func(ctx context.Context, event *types.RejectEvent) error {
		if normJSON, err := json.MarshalIndent(event, "", "  "); err == nil {
			fmt.Printf("⛔️ RejectEvent: %s\n", string(normJSON))
		}
		if event.MessageID != "" && event.ChatID != "" && event.SenderID != "" && event.Reason != "" {
			switch event.Reason {
			case "group_not_allowed":
				fmt.Printf("✅ [TC-501 Passed] 消息被策略拦截: Reason: %s, MsgID: %s\n", event.Reason, event.MessageID)
			case "dm_disabled":
				fmt.Printf("✅ [TC-502 Passed] 消息被策略拦截: Reason: %s, MsgID: %s\n", event.Reason, event.MessageID)
			case "sender_not_allowed":
				fmt.Printf("✅ [TC-503 Passed] 消息被策略拦截: Reason: %s, MsgID: %s\n", event.Reason, event.MessageID)
			case "no_mention":
				fmt.Printf("✅ [TC-504 Passed] 消息被策略拦截: Reason: %s, MsgID: %s\n", event.Reason, event.MessageID)
			case "mention_all_blocked":
				fmt.Printf("✅ [TC-505 Passed] 消息被策略拦截: Reason: %s, MsgID: %s\n", event.Reason, event.MessageID)
			default:
				fmt.Printf("✅ [TC-318 Passed] 消息被策略拦截: Reason: %s, MsgID: %s\n", event.Reason, event.MessageID)
			}
		} else {
			fmt.Printf("❌ [TC-318/5xx Failed] Reject event missing essential fields\n")
		}
		return nil
	})

	// 7. 注册卡片交互处理逻辑 (可选)
	ch.OnCardAction(func(ctx context.Context, event *types.CardActionEvent) error {
		if normJSON, err := json.MarshalIndent(event, "", "  "); err == nil {
			fmt.Printf("🃏 CardActionEvent: %s\n", string(normJSON))
		}
		operatorID := event.Operator.OpenID
		if operatorID == "" {
			operatorID = event.Operator.UserID
		}
		if operatorID != "" && event.Action.Tag != "" {
			fmt.Printf("✅ [TC-313 Passed] 收到来自 %s 的卡片交互，Tag: %s, Value: %v\n", operatorID, event.Action.Tag, event.Action.Value)
		} else {
			fmt.Printf("❌ [TC-313 Failed] Card action missing operator or action tag\n")
		}
		return nil
	})

	// 8. 注册生命周期钩子
	ch.OnReady(func() {
		fmt.Println("[WS] Client is ready and Bot Identity loaded.")
	})
	ch.OnError(func(err error) {
		fmt.Printf("[WS] Error occurred: %v\n", err)
	})
	ch.OnReconnecting(func() {
		fmt.Println("[WS] Connection lost, reconnecting...")
	})
	ch.OnReconnected(func() {
		fmt.Println("[WS] Connection re-established!")
	})
	ch.OnDisconnected(func() {
		fmt.Println("[WS] Connection disconnected.")
	})

	// 9. 启动 WebSocket 客户端
	fmt.Println("🚀 Starting Feishu Bot via Channel...")
	go func() {
		err := ch.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	select {}
}
