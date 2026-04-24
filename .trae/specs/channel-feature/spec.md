# Channel Feature Spec

## Why
当前 `oapi-sdk-go` 提供的事件接收和消息发送 API 过于底层，开发者需要手动处理复杂的事件类型转换、资源上传和富文本构建。为了提升开发体验，参考 Node SDK 的 `feature/channel` 分支，在 Go SDK 中引入高层的 `Channel` 抽象，提供极简的事件监听、消息发送以及流式响应（Streaming）能力。

## What Changes
- 新增 `channel` 包作为高阶抽象层。
- 提供 `channel.New(...)` 构造函数，封装 `lark.Client` 和可选的 `ws.Client`。
- **事件标准化 (Normalization)**：将复杂的飞书原始事件（如 `P2MessageReceiveV1`, 卡片交互事件）转换为统一的 `NormalizedMessage`, `CardActionEvent` 等易用结构。
- **增强型消息发送 (Outbound Sender)**：支持通过统一的 `Send` 方法发送文本、Markdown、图片、文件、卡片，自动处理本地文件上传和多媒体资源转换。
- **流式响应 (Streaming)**：提供 `Stream` 接口，支持对 Markdown 内容和卡片进行流式更新（自动防抖和节流）。
- **安全流水线 (Safety Pipeline)**：内置消息去重（Dedup）、并发控制（Processing Lock）机制。

## Impact
- Affected specs: 提升整体 SDK 易用性，特别是构建对话机器人的场景。
- Affected code: 新增 `channel` 目录及相关子目录（`normalize`, `outbound`, `safety` 等），不影响现有 `core`, `event`, `ws` 的基础能力。

## ADDED Requirements
### Requirement: 统一的消息发送与资源管理
The system SHALL provide 统一的 `Send` 方法：
- 自动识别并上传本地文件/内存数据，获取 `image_key` / `file_key`。
- 将 Markdown 转换为对应的飞书富文本（Post）格式或直接支持 Markdown。

#### Scenario: 自动上传并发送图片
- **WHEN** user calls `channel.Send(ctx, to, channel.ImageInput{Path: "/tmp/1.png"})`
- **THEN** SDK 会先调用飞书图片上传接口获取 image_key，然后调用消息发送接口发送该图片消息。

### Requirement: 流式 Markdown 与 卡片更新
The system SHALL provide `Stream` 方法：
- 开发者可以使用 Stream 接口高频次传入文本片段。
- SDK 底层会自动缓冲（Throttle）并以合并更新（EditMessage）的方式平滑呈现到飞书客户端。

## MODIFIED Requirements
无。此功能为纯新增高阶封装，不改变底层逻辑。

## REMOVED Requirements
无。
