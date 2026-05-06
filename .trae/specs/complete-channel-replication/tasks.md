# Tasks
- [x] Task 1: 补充边缘事件类型的 Normalization (归一化)
  - [x] SubTask 1.1: 在 `channel/types` 定义 `ReactionEvent`, `CommentEvent`, `BotAddedEvent` 数据结构。
  - [x] SubTask 1.2: 在 `channel/normalize/parser.go` 中实现 `ParseReaction`, `ParseComment`, `ParseBotAdded`。
  - [x] SubTask 1.3: 扩展 `channel.Channel` 接口及其实现，增加 `OnReaction`, `OnComment`, `OnBotAdded` 方法。

- [x] Task 2: 补充 WebSocket 客户端生命周期回调
  - [x] SubTask 2.1: 在 `ws/client.go` 中添加 `WithOnReady`, `WithOnError`, `WithOnReconnecting`, `WithOnReconnected` 等 ClientOption。
  - [x] SubTask 2.2: 在 WebSocket 连接/重连逻辑中正确触发这些回调函数。

- [x] Task 3: 补充 HTTP 客户端 User-Agent `source` 选项
  - [x] SubTask 3.1: 在 `core/client.go` (或项目相关 Client 选项) 中新增 `source` 参数，以便其在 HTTP Header 中附带 `source/<name>`。

- [x] Task 4: 实现发送侧的容错与回退机制 (Outbound Sender Fallback)
  - [x] SubTask 4.1: 在 `channel/send.go` 中实现 `sendOneWithFallback`：当富文本格式错误时退化为纯文本发送。
  - [x] SubTask 4.2: 在 `channel/send.go` 中处理 `reply_to` 目标不存在的报错，降级为普通新建消息发送。

- [x] Task 5: 补全卡片流式控制器与队列更新 (Card Stream & Update Queue)
  - [x] SubTask 5.1: 在 `channel/stream.go` (或新建 `card_stream.go`) 实现 `CardStreamController`。
  - [x] SubTask 5.2: 在 `channel/stream.go` 中实现 `UpdateQueue` 串行更新队列（解决异步高频请求被乱序执行的问题）。

- [x] Task 6: 补全项目文档与图片资源
  - [x] SubTask 6.1: 复制 Node SDK 的 `docs/channel.md` 与 `docs/channel.zh.md` 至当前项目的 `doc/channel.md` 与 `doc/channel.zh.md` 并适配为 Go 代码的语法。
  - [x] SubTask 6.2: 复制相关图片到 `doc/assets/` 目录下（例如 `debugger-tip.png`, `msg-card.png`, `deprecated.png` 等）。
  - [x] SubTask 6.3: 更新项目根目录的 `README.md`，添加关于 Channel 高级模块的简介及相关指引。

# Task Dependencies
- Task 1 无依赖
- Task 2, 3 涉及 SDK 底层 Client，需优先或并行处理
- Task 4, 5 涉及 Channel 核心的增强逻辑
- Task 6 (文档与图片) 可以在其他功能完成后执行

# 架构设计约束
1. **WS Client 生命周期回调**：在 `ws/client.go` 中使用 Functional Options 模式 (`WithOnReady(func())` 等) 实现。
2. **User-Agent `source` 追加**：在根目录 `lark.ClientOption` 增加 `WithSource(string)`，在 `core/client.go` 中统拦截并修改 Header，实现全局生效。
3. **CardStream 与串行队列**：采用 Go 原生的 Channel + 后台 Goroutine Worker 模式实现 `UpdateQueue`，天然防并发并保证串行更新，卡片数据统一接收 `interface{}` 格式。
