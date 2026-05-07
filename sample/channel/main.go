package main

import (
	"context"
	"encoding/base64"
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
	flag.StringVar(&appID, "app_id", "", "Lark App ID (也可以通过 APP_ID 环境变量传入)")
	flag.StringVar(&appSecret, "app_secret", "", "Lark App Secret (也可以通过 APP_SECRET 环境变量传入)")
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
	ch := channel.NewChannel(client, wsClient)

	// 5. 注册消息处理逻辑
	ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
		fmt.Printf("收到来自 %s 的消息: %s\n", msg.UserID, msg.Content)

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
			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Text:      "Echo: " + msg.Content + "\n\n💡 支持指令:\n/stream\n/markdown\n/card\n/cardstream\n/mention\n/file\n/image\n/post\n/sharechat\n/shareuser",
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	// 6. 注册边缘事件与生命周期
	ch.OnBotAdded(func(ctx context.Context, event *types.BotAddedEvent) error {
		fmt.Printf("[事件] 机器人被添加到群聊: %s (ChatID: %s) by UserID: %s\n", event.ChatName, event.ChatID, event.UserID)
		return nil
	})

	ch.OnReaction(func(ctx context.Context, event *types.ReactionEvent) error {
		fmt.Printf("[事件] 收到表情回应: %s, action: %s, messageID: %s\n", event.ReactionType, event.Action, event.MessageID)
		return nil
	})

	ch.OnComment(func(ctx context.Context, event *types.CommentEvent) error {
		fmt.Printf("[事件] 收到文档评论/回复: parentID: %s, commentID: %s\n", event.ParentID, event.MessageID)
		return nil
	})

	// 7. 注册卡片交互处理逻辑 (可选)
	ch.OnCardAction(func(ctx context.Context, msg *types.NormalizedMessage) error {
		fmt.Printf("收到来自 %s 的卡片交互，Value: %v\n", msg.UserID, msg.Content)
		return nil
	})

	// 8. 注册生命周期钩子
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

	// Give it a moment to connect and fetch bot identity
	time.Sleep(2 * time.Second)
	botInfo := ch.GetBotIdentity(context.Background())
	if botInfo != nil {
		displayAppID := botInfo.AppID
		if displayAppID == "" {
			displayAppID = appID
		}
		fmt.Printf("🤖 Bot Identity Loaded: AppID=%s, OpenID=%s\n", displayAppID, botInfo.OpenID)
	} else {
		fmt.Println("⚠️ Failed to load Bot Identity.")
	}

	select {}
}
