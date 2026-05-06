# Channel module

`Channel` is a high-level module built on top of `ws.Client` / `event.EventDispatcher` / `lark.Client`. It bundles the chores of running a Feishu bot — transport, message normalization, safety policy, outbound sending, streaming replies, media upload, card interactions — so you can focus on the business logic.

**When to use Channel**: conversational bots (AI chat, streaming replies, interactive card buttons, media handling, mention-all policy, etc.). If you only need to receive a few events and do simple processing, `ws.Client` + `event.EventDispatcher` is sufficient.

---

## Table of contents

- [Minimal example](#minimal-example)
- [Event listening](#event-listening)
- [`NormalizedMessage` fields](#normalizedmessage-fields)
- [Sending messages](#sending-messages)
- [Streaming replies](#streaming-replies)
- [Error types](#error-types)

---

## Minimal example

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

	// 1. Initialize API Client
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(core.LogLevelInfo))

	// 2. Initialize WebSocket Client
	wsClient := larkws.NewClient(appID, appSecret, larkws.WithLogLevel(core.LogLevelInfo))

	// 3. Create Channel
	ch := channel.NewChannel(client, wsClient)

	// 4. Register Event Handlers
	ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
		fmt.Printf("received message: %s\n", msg.Content)

		_, err := ch.Send(ctx, &types.SendInput{
			ChatID:         msg.ChatID,
			Markdown:       fmt.Sprintf("received: %s", msg.Content),
			ReplyMessageID: msg.MessageID,
		})
		return err
	})

	// 5. Start WebSocket connection
	err := wsClient.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// 6. Wait for termination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
```

---

## Event listening

```go
ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
    // text / post / media messages
    return nil
})

ch.OnReaction(func(ctx context.Context, event *types.ReactionEvent) error {
    // emoji reaction add/remove
    return nil
})

ch.OnBotAdded(func(ctx context.Context, event *types.BotAddedEvent) error {
    // bot invited into a group
    return nil
})

ch.OnComment(func(ctx context.Context, event *types.CommentEvent) error {
    // doc comment mentioning the bot
    return nil
})
```

---

## `NormalizedMessage` fields

Payload for the `message` event. Feishu's various message types (text / post / image / file / audio / video / sticker / share_chat / share_user / merge_forward / ...) are converted into a uniform markdown + XML-style-tag format, so you don't need to parse the raw Feishu message structure yourself.

| Field | Type | Description |
| --- | --- | --- |
| `EventID` | string | Unique event ID for deduplication |
| `MessageID` | string | Message id |
| `ChatID` | string | Chat id |
| `ChatType` | string | Direct message (`p2p`) / group (`group`) |
| `UserID` | string | Sender user id |
| `Content` | string | Normalized content (markdown, media rendered as XML-style tags) |
| `Mentions` | `[]types.Mention` | @ list (excludes bot itself) |
| `MentionAll` | boolean | Whether `@all` was used |
| `MentionedBot` | boolean | Whether the bot itself was mentioned |
| `Resources` | `[]types.Resource` | Resources referenced in the message (image / file / audio / video / sticker) |
| `CreateTimeMs` | int64 | Millisecond timestamp |
| `RawEvent` | interface{} | Raw event body |

---

## Sending messages

`ch.Send(ctx, input)` — supports various message inputs:

```go
// Text-based
ch.Send(ctx, &types.SendInput{ChatID: chatID, Text: "plain text"})
ch.Send(ctx, &types.SendInput{ChatID: chatID, Markdown: "hello **world**"})
ch.Send(ctx, &types.SendInput{ChatID: chatID, Card: cardJsonV2})

// Media (Uploader automatically uploads them before sending)
ch.Send(ctx, &types.SendInput{ChatID: chatID, ImagePath: "./fixtures/sample.png"})
ch.Send(ctx, &types.SendInput{ChatID: chatID, FilePath: "./fixtures/doc.pdf"})

// Reply to a message
ch.Send(ctx, &types.SendInput{
    ChatID: chatID, 
    Markdown: "please check",
    ReplyMessageID: msg.MessageID,
    Mentions: msg.Mentions, // structured @ mentions
})
```

**Do not hand-write `@username` in the text.** Pass `[]types.Mention` via `Mentions` and the SDK will assemble the Feishu @ placeholders correctly.

### Automatic fallbacks

Outbound sending has built-in fallbacks that are transparent to the caller:

- **Reply target revoked** (`target_revoked`) → drop `ReplyMessageID` and resend as a fresh message
- **Post schema rejected** (`format_error`) → fall back to a plain-text resend

If the fallback also fails, an error is returned.

---

## Streaming replies

Keep the call site non-streaming; the SDK owns throttling and the native typewriter animation (via cardkit v1):

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

- A `"Thinking..."` placeholder card is sent before the first `Append`.
- The server renders the typewriter animation (no client-side throttling needed).
- If the producer never produces anything → card shows `"(no content)"`.
- Call `Close(ctx)` to complete the stream and finalize the card.

---

## Error types

All outbound failures return specific error wrappers, allowing you to distinguish error scenarios like:

- `types.ErrFormatError`: Message schema validation failed
- `types.ErrTargetRevoked`: Reply target deleted / recalled
- `types.ErrPermissionDenied`: Auth failure / insufficient scope
- `types.ErrRateLimited`: Rate limiting (HTTP 429)
- `types.ErrUploadFailed`: Media upload failed

---

## Common issues

1. **Don't write `@username` as literal text** in replies — use the structured `SendInput.Mentions`.
2. **Group messages require @bot by default.** To receive all group messages, request `im:message.group_msg` (admin approval required).
3. **Card buttons not firing** is usually either a missing `card.action.trigger` subscription or a v1 card schema. V2 uses `column_set` → `column` → `button` with `behaviors: [{ type: 'callback', value }]`.
4. **Use `Stream()`** for streaming output; don't loop `Send()` / update card by hand.
5. **WebSocket client** waits for the real handshake; runtime disconnects are auto-retried with corresponding events emitted.
