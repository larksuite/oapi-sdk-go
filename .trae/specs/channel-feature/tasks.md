# Tasks

- [x] Task 1: 初始化 `channel` 包结构与核心类型定义
  - [x] SubTask 1.1: 创建 `channel` 目录及 `types.go`，定义 `Channel` 接口、`NormalizedMessage`、`SendInput` 等基础结构。
  - [x] SubTask 1.2: 实现 `Channel` 构造函数，支持注入 `lark.Client` 和 `ws.Client`。

- [x] Task 2: 实现事件标准化处理 (Normalization)
  - [x] SubTask 2.1: 实现 `Message` 事件的标准化解析，提取文本、资源(图片/文件)、Mentions。
  - [x] SubTask 2.2: 实现 `CardAction` 等其他常见事件的标准化。
  - [x] SubTask 2.3: 在 `Channel` 提供 `OnMessage`, `OnCardAction` 等快捷注册方法。

- [x] Task 3: 实现增强型消息发送 (Outbound Sender)
  - [x] SubTask 3.1: 实现资源上传器 (Uploader)，封装图片、音频、文件的自动上传。
  - [x] SubTask 3.2: 实现 Markdown 转换器，将标准 Markdown 转为飞书的富文本 (Post) 格式。
  - [x] SubTask 3.3: 实现 `channel.Send` 主逻辑，支持按类型路由和发送消息。

- [x] Task 4: 实现流式响应能力 (Streaming)
  - [x] SubTask 4.1: 实现 Throttle 控制器（节流与合并）。
  - [x] SubTask 4.2: 实现 `MarkdownStreamController`，支持 `Append` 文本并通过 Edit 接口更新消息。
  - [x] SubTask 4.3: 实现 `channel.Stream` 入口。

- [x] Task 5: 实现安全流水线 (Safety Pipeline)
  - [x] SubTask 5.1: 实现 LRU Dedup Cache 处理重复事件。

- [x] Task 6: 编写单元测试与示例
  - [x] SubTask 6.1: 补充 `channel` 核心方法的单元测试。
  - [x] SubTask 6.2: 在 `sample/` 目录下提供使用 `channel` 的完整对话机器人示例。

# Task Dependencies
- Task 2 depends on Task 1
- Task 3 depends on Task 1
- Task 4 depends on Task 3
- Task 6 depends on all previous Tasks
