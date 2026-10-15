# Channel Sample Testing Guide

本目录提供两类测试入口：

- `sample/channel/main.go`：面向真实飞书环境的手工联调入口，适合验证事件接收、策略拦截、卡片交互、流式更新等场景。
- `sample/channel_test_cases/main.go`：面向真实飞书环境的自动化/半自动回归入口，适合批量验证发送、更新、回退、下载与策略能力。

## 1. 前置条件

在开始前，请确认以下条件已经满足：

1. 已创建并启用一个飞书应用，并拿到 `APP_ID` 与 `APP_SECRET`。
2. 已在开发者后台为应用开启长连接事件订阅。
3. 已为应用开通消息发送、消息接收、文件上传/下载、名片分享等相关权限。
4. 机器人已被加入待测试群聊，或你已准备一个可接收消息的单聊用户。
5. 如需测试邮箱转用户、名片分享等能力，应用需具备对应联系人读取权限。

建议优先使用环境变量：

```bash
export APP_ID=cli_xxx
export APP_SECRET=xxx
```

自动化测试还需要以下其一：

```bash
export RECEIVE_ID=ou_xxx_or_oc_xxx
```

或：

```bash
export EMAIL=someone@example.com
```

## 2. 手工联调入口

启动命令：

```bash
go run sample/channel/main.go
```

也可以显式传参：

```bash
go run sample/channel/main.go \
  -app_id="$APP_ID" \
  -app_secret="$APP_SECRET" \
  -dm_mode=open \
  -respond_all=true
```

支持参数：

- `-app_id` / `-app_secret`：飞书应用凭证。
- `-dm_mode`：单聊策略，支持 `open`、`disabled`、`allowlist`。
- `-respond_all`：群聊中是否响应 `@all`。

启动成功后，程序会持续监听消息、表情、评论、卡片交互、策略拒绝事件，并在终端输出 `TC-xxx Passed/Failed` 日志。

## 3. 手工测试建议顺序

建议按下面顺序逐步验证：

1. 先发送一条普通文本，确认机器人已成功接收消息并能回显。
2. 再测试图片、文件、富文本、卡片、流式更新。
3. 最后测试策略拦截、表情事件、文档评论和卡片交互等边缘能力。

## 4. 手工测试指令

向机器人发送以下指令，可以触发对应能力：

| 指令                                               | 目的                   | 预期                         |
|--------------------------------------------------|----------------------|----------------------------|
| 普通文本                                             | 验证 `OnMessage` 与默认回显 | 输出 `TC-301` / `TC-312` 等日志 |
| `/policy`                                        | 查看当前策略               | 机器人回复当前策略 JSON             |
| `/policy dm_mode=disabled`                       | 禁用单聊                 | 后续单聊消息触发 `TC-502`          |
| `/policy dm_mode=open`                           | 恢复单聊放行               | 单聊消息恢复正常处理                 |
| `/policy dm_mode=allowlist, dm_allowlist=ou_xxx\|ou_yyy`               | 配置单聊白名单                    | 白名单内用户放行，其他用户触发 `TC-503` |
| `/policy respond_all=false`                      | 禁止响应 `@all`          | 群聊 `@all` 触发 `TC-505`      |
| `/policy respond_all=true`                       | 允许响应 `@all`          | 群聊 `@all` 恢复正常处理           |
| `/policy require_mention=true`                   | 群聊要求显式 @机器人          | 未提及时触发 `TC-504`            |
| `/policy require_mention=false`                  | 群聊不强制 @机器人           | 普通群消息可直接处理                 |
| `/policy group_allowlist=oc_xxx\|oc_yyy`                              | 配置群白名单                     | 非白名单群消息触发 `TC-501` |
| `/policy group_allowlist=empty`                  | 清空群白名单               | 群白名单限制被移除                  |
| `/policy dm_allowlist=empty`                     | 清空单聊白名单              | 便于恢复或切换策略测试                |
| `/stream`                                        | 测试 Markdown 流式输出     | 看到逐步更新的消息，终端可观察流式日志        |
| `/markdown`                                      | 测试 Markdown 发送       | 成功发送富文本 Markdown           |
| `/card`                                          | 测试交互卡片发送             | 发送一张带按钮的卡片                 |
| `/cardstream`                                    | 测试卡片流式更新             | 卡片内容逐步更新进度                 |
| `/mention`                                       | 测试 Mentions 组合发送     | 自动 @ 当前发送者                 |
| `/file`                                          | 测试本地临时文件上传发送         | 收到文件消息                     |
| `/image`                                         | 测试本地临时图片上传发送         | 收到图片消息                     |
| `/post`                                          | 测试直接发送 Post JSON     | 收到富文本消息                    |
| `/sharechat`                                     | 测试分享当前群名片            | 收到当前群的群名片或明确报错             |
| `/sharechat oc_d03d81dbd26b9b37b1426e6491baa095` | 测试分享指定群名片            | 收到指定群的群名片或明确报错             |
| `/shareuser someone@example.com`                 | 测试个人名片发送             | 查到用户后发送个人名片                |

策略指令支持逗号分隔多个参数，例如：

```text
/policy dm_mode=allowlist, dm_allowlist=ou_xxx|ou_yyy
/policy group_allowlist=oc_group_1|oc_group_2, require_mention=true
```

## 4.1 策略手工联调步骤

如果你想专门验证 `TC-501 ~ TC-505`，建议按下面顺序操作：

1. 先发送 `/policy`，确认当前策略。
2. 根据目标场景切换策略。
3. 再发送一条能命中该策略的真实消息。
4. 观察终端中的 `RejectEvent` 与 `TC-50x Passed` 输出。

常见场景如下：

| TC       | 操作                                                                      | 预期                       |
|----------|-------------------------------------------------------------------------|--------------------------|
| `TC-501` | `/policy group_allowlist=oc_allowed_group` 后，在其他群发送消息                   | 输出 `group_not_allowed`   |
| `TC-502` | `/policy dm_mode=disabled` 后，从单聊发送消息                                    | 输出 `dm_disabled`         |
| `TC-503` | `/policy dm_mode=allowlist, dm_allowlist=ou_allowed_user` 后，用其他用户单聊发送消息 | 输出 `sender_not_allowed`  |
| `TC-504` | `/policy require_mention=true` 后，在群里发送未 @机器人的消息                         | 输出 `no_mention`          |
| `TC-505` | `/policy respond_all=false, require_mention=false` 后，在群里发送 `@all` 消息    | 输出 `mention_all_blocked` |

建议每轮策略测试结束后恢复到较宽松的状态，例如：

```text
/policy dm_mode=open, require_mention=false, respond_all=true, group_allowlist=empty, dm_allowlist=empty
```

## 5. 手工 TC 对照

`sample/channel/main.go` 中已经内置了下列手工验证编号：

### 消息与基础事件

- `TC-301`：接收文本消息
- `TC-302`：接收图片消息
- `TC-303`：接收 Post 消息
- `TC-304`：接收文件消息
- `TC-305`：接收音频消息
- `TC-306`：接收视频消息
- `TC-307`：接收分享名片
- `TC-308`：接收交互卡片
- `TC-309`：接收合并转发消息
- `TC-310`：识别提及机器人
- `TC-311`：识别 `@all`
- `TC-312`：识别单聊/群聊类型

### 交互与边缘事件

- `TC-313`：卡片交互回调
- `TC-314`：表情新增事件
- `TC-315`：表情移除事件
- `TC-316`：机器人被添加到群聊
- `TC-317`：文档评论/回复事件
- `TC-318`：策略拒绝事件总入口

### 策略拦截

- `TC-501`：群聊白名单拦截
- `TC-502`：单聊禁用拦截
- `TC-503`：单聊白名单拦截
- `TC-504`：群聊未提及机器人拦截
- `TC-505`：`@all` 拦截

## 6. 自动化测试入口

启动全部自动化/半自动用例：

```bash
go run sample/channel_test_cases/main.go
```

显式传参方式：

```bash
go run sample/channel_test_cases/main.go \
  -app_id="$APP_ID" \
  -app_secret="$APP_SECRET" \
  -receive_id="$RECEIVE_ID"
```

如果没有 `receive_id`，可以通过邮箱查找：

```bash
go run sample/channel_test_cases/main.go \
  -app_id="$APP_ID" \
  -app_secret="$APP_SECRET" \
  -email="$EMAIL"
```

只运行某个前缀或单个用例：

```bash
go run sample/channel_test_cases/main.go -case=TC-001
go run sample/channel_test_cases/main.go -case=TC-10
go run sample/channel_test_cases/main.go -case=TC-70
```

说明：

- `-case=TC-001`：只跑单个用例。
- `-case=TC-10`：会跑 `TC-101` 到 `TC-109`、`TC-110` 等同前缀用例。
- `-case=TC-70`：会跑重试与发送回退相关用例。

## 7. 自动化覆盖范围

当前 `sample/channel_test_cases/main.go` 已覆盖以下大类：

- `TC-001` ~ `TC-007`：连接、身份、错误监听、重连监听
- `TC-101` ~ `TC-113`：文本、Markdown、长文本分片、Post、图片、文件、音频、视频、分享、卡片、提及
- `TC-201` ~ `TC-208`：回复、更新、撤回、添加表情
- `TC-401` ~ `TC-403`：Markdown/Card 流式更新与异常处理
- `TC-601` ~ `TC-605`：上传/下载与 SSRF Guard
- `TC-701` ~ `TC-703`：重试与发送回退

策略相关能力建议通过 `sample/channel/main.go` 的 `/policy` 指令和 `OnReject`
日志做手工联调，并通过 `channel/safety/policy_gate_test.go`、`channel/policy_reject_test.go` 做单元测试覆盖。

## 8. 测试资源

自动化测试依赖以下媒体文件：

- `sample/channel_test_cases/test_audio.mp3`
- `sample/channel_test_cases/test_video.mp4`

如果文件不存在，脚本会回退为占位字节流，但 `TC-107`、`TC-108` 的后端验证可能不稳定。

## 9. 常见问题

### 收不到事件

请确认：

1. 飞书后台已开启长连接订阅。
2. 机器人已加入目标群聊。
3. 应用具备相应事件订阅权限。
4. 启动日志中已经看到 WebSocket ready 相关输出。

### 单聊/群聊消息没有响应

优先检查当前策略：

```text
/policy
```

如果策略限制了单聊、`@all` 或未提及机器人，则可能触发 `TC-502`、`TC-504`、`TC-505`。

### 卡片按钮点击后没有回调

请确认：

1. 已发送的是交互卡片而不是普通消息。
2. 长连接事件订阅正常。
3. 终端里是否有 `CardActionEvent` 原始输出。

### 文档评论事件没有触发

请确认应用是否已订阅对应文档评论事件，并且评论操作发生在应用可见的文档资源上。

## 10. 推荐使用方式

建议开发时采用下面的节奏：

1. 先运行 `sample/channel/main.go` 做手工联调。
2. 用文本、卡片、流式、策略相关指令做基础验证。
3. 再运行 `sample/channel_test_cases/main.go` 做回归检查。
4. 最后根据新增能力补充对应 `TC-xxx` 编号与本文档说明。
