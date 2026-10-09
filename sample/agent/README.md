# Agent v1 完整调用示例（Go）

从创建会话、发送问题到逐事件读取回复，展示上传、游标续读和业务打断。单接口生成示例仍保留在原目录；这里是人工维护的组合调用示例。

## 前提与安装

Go 1.17 或更高版本。`sample/go.mod` 的 replace 已指向当前仓库，直接使用本地 SDK，无需另建模块。

本示例要求 SDK 源码包含 Agent v1、SSE 及文件名覆盖能力。当前说明不承诺某个已经发布的包版本包含这些能力；先使用当前 checkout，发布后再切换到实际包含这些能力的版本。

应用须开通 Agent 能力，用户能够正常使用 Agent，并提供有效的用户 access token。应用接口权限及用户授权须覆盖：

| 操作 | Scope |
|---|---|
| 创建会话 | `doubao:session:write` |
| 发送消息/业务打断 | `doubao:session.event:write` |
| SSE 读取 | `doubao:session.event:read` |
| 上传文件 | `doubao:file:write` |
| 打断后查询状态 | `doubao:session.message:read` |

示例从环境变量读取凭证，不自动读取 `.env`。在本机设置环境变量，勿将真实凭证写入示例或提交到 Git：

```bash
export APP_ID='你的应用ID'
export APP_SECRET='你的应用Secret'
export USER_ACCESS_TOKEN='有效的用户AccessToken'
export OPEN_BASE_URL='https://open.feishu.cn'
```

## 1. 完整文本对话

```bash
cd sample
go run ./agent chat
```

流程：创建空会话 → 发送 `user.message` → 取得返回的用户 query 消息 ID → 用 session ID 与 query ID 打开 SSE → 每条完整事件到达立即处理 → 关闭流。

输出会话 ID、query ID、assistant ID、事件类型/文本长度及事件游标。默认不输出完整用户内容。查看源码中 `handleEvent` / `handle_event`，将解析出的 `text` 或 `delta` 接入自己的 UI。

请求中的文本内容必须显式设置 `type="text"`。SSE 的 `data` 是 JSON 字符串，需应用自行解析。`text.delta` 是预览增量，完整 `text` 是最终 item；应用可按 item ID 用完整内容替换预览，避免重复累加。

## 2. 文件上传与在消息中引用

默认上传示例内存文本内容，实际文件段与表单字段均使用 `report.txt`：

```bash
go run ./agent upload
```

上传真实文件（相对路径以运行命令的当前目录为准）：

```bash
export AGENT_FILE_PATH='/absolute/path/to/local.txt'
export AGENT_FILENAME='开发者报告.txt'
go run ./agent upload
```

显式 `AGENT_FILENAME` 决定实际上传名，不修改磁盘文件名。文件名应带匹配内容的扩展名，不能包含目录、换行或 NUL。读取输入文件后由示例关闭资源；Java 内存示例创建的临时文件会删除。未设置覆盖名称的通用 SDK 规则仍可从文件对象推断名称。

上传成功输出 `AGENT_FILE_URI`。将这个值原样设置，再运行 chat；示例会在 contents 中添加文件引用：

```bash
export AGENT_FILE_URI='上一步返回的URI'
go run ./agent chat
unset AGENT_FILE_URI
```

文件输入的原始对象为：

```json
{"type":"file","uri":"上传接口返回的URI"}
```

URI 不要拼接域名，不要 Base64 编码。文件 URI 使用与当前 UAT 相同用户/应用可访问的资源。当前 Meta 未列出这个请求 union 分支，示例沿用 SDK 的轻量原始对象构造能力；不增加自动类型推断或 SDK 业务校验。是否能处理文件由 Agent 服务和权限决定。

## 3. 游标续读

从 chat 输出中取得下面三个值。URL 中的 query ID 是用户消息 ID；SSE id 可能以 assistant 消息 ID 开头，二者不能互换：

```bash
export AGENT_SESSION_ID='chat 输出的会话ID'
export AGENT_QUERY_MESSAGE_ID='chat 输出的用户query消息ID'
export AGENT_LAST_EVENT_ID='成功处理过的事件的完整SSE id'
go run ./agent resume
```

示例将游标作为 `Last-Event-ID` 请求头发送，原样使用服务端值，不自行拼接 `message_id:seq`。只有成功处理事件后才记录游标；空 id 表示重置，生产应用应删除已存的检查点。

SDK 不自动重连、回放去重或聚合文本。服务端只回放其持久化事件，预览 delta 不保证可回放；断线后处理完整 item 可以恢复内容。续读返回哪些事件取决于服务端日志，客户端无法补回缺失的服务端事件。

## 4. 关闭接收连接

设置上述 session/query ID 后，执行：

```bash
go run ./agent close
```

接收一条完整事件后立即退出读取并关闭 HTTP 连接。这仅停止本客户端接收，Agent 仍可能继续生成。普通 EOF 表示此次流结束；网络错误、超时与业务失败应分别处理。

## 5. 业务打断与查询状态

在另一终端、回复还在生成时，用 chat 输出的 assistant ID 发起打断：

```bash
export AGENT_SESSION_ID='会话ID'
export AGENT_ASSISTANT_MESSAGE_ID='正在生成的assistant消息ID'
go run ./agent interrupt
```

请求事件为 `{"type":"user.interrupt","message_id":"assistant消息ID"}`。它操作业务回复，与关闭 SSE 不同。示例随后调用 Messages Get 查询一次状态；即时结果仍可能为 `in_progress`，并不代表打断失败。`cancelled` 表示被取消，`completed` 可能表示发送打断前已完成；需要等待终态时由应用进行有界查询。

## 超时、错误与幂等

普通 HTTP 请求超时为 10 秒；SSE 打开 10 秒、空闲读取 60 秒、总时长 180 秒、单事件最大 1 MiB，可在源码的 options 配置处调整。普通 API 返回非零 code、SSE 数据中的业务错误、流读取异常都会结束示例并以非零状态退出。HTTP 200 并不必然代表 SSE 业务成功。

每次独立写操作生成 UUID。示例不自动重试；若应用重试同一次写操作，应保留同一个请求体和 UUID，避免重复创建会话/消息。资源会保留在服务端，示例不自动归档会话。

## 源码导航

- `main.go`：环境变量、命令与 Context 取消。
- `conversation.go`：创建会话、轻量 union、JSON 事件处理。
- `stream.go`：逐事件读取、关闭、游标、业务打断。
- `upload.go`：Reader/真实文件上传。

Go 使用 `Next()` 取得下一条完整事件，`Event()` 读取当前事件，循环后 `Err()` 检查读取错误；defer 关闭流。整个示例支持 Context 取消，Ctrl+C 取消接收，不发送 user.interrupt。
