# Complete Channel Replication Spec

## Why
Node SDK 的 `feature/channel` 分支包含了极其丰富的对话机器人抽象能力。当前 Go SDK 项目分支中虽然已经完成了核心骨架（解析、发送、防抖流式、重试去重等）的复刻，但在边缘事件类型、发送回退容错机制、卡片流式更新、长连接生命周期钩子，以及官方文档与示例配图方面仍有遗漏。
为了实现 **100% 完整对齐**，我们需要补齐这些缺失的能力与文档，确保开发者在 Go SDK 中获得与 Node SDK 完全一致的高级开发体验。

## What Changes
- **补全事件解析 (Event Normalization)**：新增对 `BotAddedEvent`（机器人被添加）、`ReactionEvent`（表情回复）、`CommentEvent`（文档评论）的标准化解析与通道监听器。
- **补全发送容错 (Outbound Fallback)**：在 `Send` 中增加发送失败回退机制（例如富文本发送失败降级为纯文本，回复目标不存在时降级为新建消息）。
- **补全高级流式 (Advanced Streaming)**：新增 `CardStreamController` 与 `UpdateQueue`，支持卡片内容的流式/排队更新。
- **补全客户端生命周期与追踪**：
  - 在 WebSocket 客户端 `ws.Client` 中补充生命周期回调函数（`OnReady`, `OnError`, `OnReconnecting`, `OnReconnected`）。
  - 在 HTTP 客户端中补充 `source` 选项，用于自定义 `User-Agent` 追踪请求来源。
- **补全项目文档与资源 (Documentation & Assets)**：
  - 引入完整的 `channel.md` 和 `channel.zh.md` 到 `doc/` 目录。
  - 将 Node SDK 中的说明配图（`debugger-tip.png`, `msg-card.png` 等）迁移到本项目的 `doc/assets/` 下。
  - 修改 `README.md`，增加 Channel 模块的使用说明。

## Impact
- Affected specs: `channel.Channel` 接口增加新的监听器，`ws.Client` 增加新选项。
- Affected code: `channel/normalize/parser.go`, `channel/send.go`, `channel/stream.go`, `ws/client.go`, `client.go`, 根目录 `README.md` 及 `doc/`。

## ADDED Requirements
### Requirement: 边缘事件处理
系统必须支持通过 `channel.OnBotAdded`, `channel.OnReaction`, `channel.OnComment` 等方法快捷监听并解析对应事件。

### Requirement: 发送失败回退机制
当调用 `channel.Send` 时：
- **WHEN** 发送富文本（Post）由于格式问题被飞书 API 拒绝。
- **THEN** 自动将其退级为纯文本（Text）并重新尝试发送。
- **WHEN** 回复某条消息时，目标消息已被撤回或删除。
- **THEN** 自动退级为普通发送（不带 `replyTo`）。

### Requirement: WebSocket 生命周期钩子
在长连接断开、重连、连接成功时，允许开发者通过注册的 Callback 函数获取通知。

### Requirement: 完善的开发文档
包含英/中双语的 Channel 教程文档与 README 简介，并带有相关解释配图。
