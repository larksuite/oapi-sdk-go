# OAPI-SDK-GO

## New SDK version v2 (SDK v2 新版本)

> ### 使用飞书（[feishu.cn](http://open.feishu.cn)），请看文档：[FeiShu.md](FeiShu.md)
> ### SDK v2 advantage（SDK v2 优势）
> - SDK V2 introduces fewer packages, avoiding the problem of package name conflict between SDK and developer`s business system
> （SDK v2 引入的包更少，避免 SDK 包名与开发者的业务系统出现包名冲突的问题)
>
>
> - SDK V2 request API need parameters are encapsulated into a structure (including HTTP request path, Query, and body parameters) to avoid the omission caused by scattered Settings of parameters required by the request
> （SDK v2 请求 API 需要的参数，封装成一个结构体（含：HTTP request path、query、body 参数），避免请求需要的参数分散设置，出现遗漏）
>
>
> - The SDK V2 adds models for various messages to facilitate building message content
> （SDK v2 增加了各种消息内容 Model，方便构建消息内容）
> 
> - The SDK V2 adds sending messages via custom bots
> （SDK v2 增加了通过自定义机器人发送消息）


## Older SDK version v1 (SDK v1 老版本）

> SDK V1 source code is not in the main branch, please cut to the specific tag branch to see（SDK v1 的源代码不在 main
  分支，可以切到具体的 tag 版本上看）
>
>
> - 使用飞书（[feishu.cn](http://open.feishu.cn)），请看：[doc/FeiShu.old.md](doc/FeiShu.old.md)
> （SDK v1 与 v2 可以在 Go module 项目中同时使用，也可以使用 v2 版本进行重构）
>
>
> - Use Lark([larksuite.com](http://open.larksuite.com)), see: [doc/LarkSuite.old.md](doc/LarkSuite.old.md)
> （SDK V1 and V2 can be used in Go Module project at the same time, or v2 version can be used for reconstruction）
>
