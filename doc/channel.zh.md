# Channel 模块

`Channel` 是在 `ws.Client` / `event.EventDispatcher` / `lark.Client` 之上封装的**高层模块**，把飞书机器人接入过程中的传输、消息归一化、安全策略、出站发送、流式回复、媒体上传、卡片交互等杂活封装好，使用者只需要专注业务逻辑。

**什么时候用 Channel**：需要做会话式机器人（AI 对话、流式回复、卡片按钮、媒体上传、@所有人策略等）时。若只是收少量事件做简单处理，`ws.Client` + `event.EventDispatcher` 组合已经够用。

---

## 目录

- [最小示例](#最小示例)
- [事件监听](#事件监听)
- [`NormalizedMessage` 字段](#normalizedmessage-字段)
- [发送消息](#发送消息)
- [流式回复](#流式回复)
- [错误类型](#错误类型)

---

## 最小示例

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/channel"
	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

func main() {
	appID := os.Getenv("FEISHU_APP_ID")
	appSecret := os.Getenv("FEISHU_APP_SECRET")

	// 1. 初始化 API Client
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(core.LogLevelInfo))

	// 2. 初始化 WebSocket Client
	wsClient := larkws.NewClient(appID, appSecret, larkws.WithLogLevel(core.LogLevelInfo))

	// 3. 创建 Channel
	ch := channel.NewChannel(client, wsClient)

	// 4. 注册事件处理器
	ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
		fmt.Printf("收到消息：%s\n", msg.Content)

		_, err := ch.Send(ctx, &types.SendInput{
			ChatID:         msg.ChatID,
			Markdown:       fmt.Sprintf("已收到：%s", msg.Content),
			ReplyMessageID: msg.MessageID,
		})
		return err
	})

	// 5. 启动 WebSocket 连接
	err := wsClient.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// 6. 阻塞等待退出信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
```

---

## 事件监听

```go
ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
    // 文本 / 富文本 / 媒体消息
    return nil
})

ch.OnReaction(func(ctx context.Context, event *types.ReactionEvent) error {
    // emoji 表态增减
    return nil
})

ch.OnBotAdded(func(ctx context.Context, event *types.BotAddedEvent) error {
    // bot 被拉进群
    return nil
})

ch.OnComment(func(ctx context.Context, event *types.CommentEvent) error {
    // 文档内 @bot 评论
    return nil
})
```

---

## `NormalizedMessage` 字段

`message` 事件的 payload。飞书各种消息类型（text / post / image / file / audio / video / sticker / share_chat / share_user / merge_forward 等）都会被转换成统一的 markdown + XML-style 标签格式，开发者不需要自己解析飞书原生消息结构。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `EventID` | string | 用于去重的唯一事件 ID |
| `MessageID` | string | 消息 ID |
| `ChatID` | string | 会话 ID |
| `ChatType` | string | 单聊 (`p2p`) / 群聊 (`group`) |
| `UserID` | string | 发送者用户 ID |
| `Content` | string | 归一化后的文本内容（markdown 格式，媒体用 XML-style 标签表示） |
| `Mentions` | `[]types.Mention` | @ 列表（不含 bot 自己） |
| `MentionAll` | boolean | 是否 @所有人 |
| `MentionedBot` | boolean | 是否 @了 bot 本身 |
| `Resources` | `[]types.Resource` | 消息里引用的资源（image / file / audio / video / sticker） |
| `CreateTimeMs` | int64 | 毫秒时间戳 |
| `RawEvent` | interface{} | 原始事件体 |

---

## 发送消息

`ch.Send(ctx, input)` —— 支持多种 input 类型：

```go
// 基础文本类
ch.Send(ctx, &types.SendInput{ChatID: chatID, Text: "plain text"})
ch.Send(ctx, &types.SendInput{ChatID: chatID, Markdown: "hello **world**"})
ch.Send(ctx, &types.SendInput{ChatID: chatID, Card: cardJsonV2})

// 媒体（传入本地路径，Uploader 会自动上传后再发送）
ch.Send(ctx, &types.SendInput{ChatID: chatID, ImagePath: "./fixtures/sample.png"})
ch.Send(ctx, &types.SendInput{ChatID: chatID, FilePath: "./fixtures/doc.pdf"})

// 回复消息
ch.Send(ctx, &types.SendInput{
    ChatID: chatID, 
    Markdown: "please check",
    ReplyMessageID: msg.MessageID,
    Mentions: msg.Mentions, // 结构化 @
})
```

**不要在文本里手写 `@用户名`**——请用 `Mentions` 结构化传入 `[]types.Mention`，SDK 负责正确拼装 Feishu 的 @ 占位符。

### 自动降级

出站发送内置自动降级：

- **回复目标已撤销**（`target_revoked`）→ 自动去掉 `ReplyMessageID`，重发为普通消息
- **post 结构校验失败**（`format_error`）→ 自动降级为纯文本重发

这两条对调用方透明。如果降级后仍然失败，向调用方返回错误。

---

## 流式回复

调用方保持非流式的写法，由 SDK 接管节流和原生打字机动画（基于 cardkit v1）：

```go
ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
    streamCtrl, err := ch.Stream(ctx, &types.SendInput{
        ChatID:         msg.ChatID,
        ReplyMessageID: msg.MessageID,
    })
    if err != nil {
        return err
    }

    for chunk := range llmStream(msg.Content) {
        streamCtrl.Append(ctx, chunk)
    }
    
    streamCtrl.Close(ctx)
    return nil
})
```

- 首次 `Append` 前自动发 "Thinking..." 占位卡片
- 服务端渲染打字机动画（无需客户端自己节流）
- 未产出任何内容 → 卡片显示 "(no content)"
- 结束时调用 `Close(ctx)` 完成流式并定格卡片状态。

---

## 错误类型

所有出站调用失败都会返回特定的错误包装类型，方便业务区分：

- `types.ErrFormatError`：消息格式校验失败
- `types.ErrTargetRevoked`：回复目标已删除 / 撤回
- `types.ErrPermissionDenied`：鉴权失败 / 权限不足
- `types.ErrRateLimited`：限流（HTTP 429）
- `types.ErrUploadFailed`：媒体上传失败

---

## 常见问题

1. **不要在回复文本里写 `@用户名`**——用 `SendInput.Mentions` 结构化传。
2. **群消息默认需要 @bot** 才触发；想响应所有群消息需要申请 `im:message.group_msg` 权限。
3. **卡片按钮不响应**通常是：没加 `card.action.trigger` 订阅 / 卡片用了 V1 schema（需要 V2：`column_set` → `column` → `button` 的 `behaviors: [{ type: 'callback', value }]` 结构）。
4. **流式用 `Stream()`**，不要手动反复 `Send()` 更新卡片模拟。
5. **WebSocket 客户端**会自动重连并等待真实握手。
