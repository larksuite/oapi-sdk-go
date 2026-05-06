* [x] `channel.Channel` 支持 `OnBotAdded`, `OnReaction`, `OnComment` 等边缘事件解析与监听。

* [x] `ws.Client` 支持在长连接状态变更时触发相应的闭包回调 (`OnReady`, `OnError` 等)。

* [x] `lark.Client` (或底层 HTTP 客户端) 支持 `WithSource("my-app")` 并在请求头 `User-Agent` 中附加 `source/my-app`。

* [x] `channel.Send` 具备在发生错误（富文本非法、被回复的消息不存在）时自动降级发送的 Fallback 机制。

* [x] `channel.Stream` / `CardStreamController` 能够对卡片和文本进行基于 `UpdateQueue` 串行排队的防抖更新。

* [x] `doc/channel.md` 与 `doc/channel.zh.md` 教程文档已创建并使用 Go 代码示例进行适配。

* [x] `doc/assets/` 目录下已包含所需截图（如 `debugger-tip.png` 等）。

* [x] `README.md` 首页已包含 Channel 模块的中文与英文简介段落。

