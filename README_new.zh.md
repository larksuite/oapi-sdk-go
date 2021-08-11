[**README of Larksuite(Overseas)**](README.md) | 飞书

# 飞书开放接口SDK

## 概述

---

- 飞书开放平台，便于企业应用与飞书集成，让协同与管理更加高效，[概述](https://open.feishu.cn/document/uQjL04CN/ucDOz4yN4MjL3gzM)

- 使用 SDK，便捷调用服务端 API 与订阅服务端 Event，例如：消息&群组、通讯录、日历、视频会议、云文档、 OKR等具体可以访问 [飞书开放平台文档](https://open.feishu.cn/document/) 看看【服务端
  API】

## 问题反馈

---

如有任何 SDK 使用相关问题，请提交 [Github Issues](https://github.com/larksuite/oapi-sdk-go/issues), 我们会在收到 Issues 的第一时间处理，并尽快给您答复

## 运行环境

---

- Golang 1.9及以上

## 安装方法

---

```shell
go get github.com/larksuite/oapi-sdk-go@v1.1.41
```

## 术语解释
- 飞书（FeiShu）：Lark 在中国的称呼，主要为国内的企业提供服务，拥有独立的[域名地址](https://www.feishu.cn)
- LarkSuite：Lark 在海外的称呼，主要为海外的企业提供服务，拥有独立的[域名地址](https://www.larksuite.com/) 
- 开发文档：开放平台的开放接口的参考，**开发者必看，可以使用搜索功能，高效的查询文档**。[更多介绍说明](https://open.feishu.cn/document/) 
- 开发者后台：开发者开发应用的管理后台，[更多介绍说明](https://open.feishu.cn/app/) 
- 企业自建应用：应用仅仅可在本企业内安装使用，[更多介绍说明](https://open.feishu.cn/document/uQjL04CN/ukzM04SOzQjL5MDN) 
- 应用商店应用：应用会在 [应用目录](https://app.feishu.cn/?lang=zh-CN) 展示，各个企业可以选择安装，[更多介绍说明](https://open.feishu.cn/document/uQjL04CN/ugTO5UjL4kTO14CO5kTN) 
  
![App type](assert/app_type.zh.png)

## 快速使用

---

### 调用服务端API

- **必看** [如何调用服务端API](https://open.feishu.cn/document/ukTMukTMukTM/uYTM5UjL2ETO14iNxkTN/guide-to-use-server-api)
  ，了解调用服务端API的过程及注意事项
  - 由于 SDK 已经封装了 app_access_token、tenant_access_token 的获取，所以在调业务API的时候，不需要去获取 app_access_token、tenant_access_token
  - 如果 API 需要使用 user_access_token，需要进行通过 lark.SetUserAccessToken("UserAccessToken") 设置，具体请看：README.zh.md -> 如何构建请求（APIRequest）
- 更多示例，请看：[sample/api/api.go](sample/api/api.go)（含：文件的上传与下载）

#### 使用`企业自建应用`访问 [发送消息API](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/create) 示例

- 在 [service](./service) 下的业务 API，都是可以直接使用SDK

```go
package demo

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	im "github.com/larksuite/oapi-sdk-go/service/im/v1"
)

var appConf *lark.AppConfig

func init() {
	// 企业自建应用的配置
	// lark.DomainFeiShu：表示访问飞书的开放 API
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// EncryptKey、VerificationToken："开发者后台" -> "事件订阅" -> 事件订阅（Encrypt Key、Verification Token）
	// HelpDeskID、HelpDeskToken, 服务台 token：https://open.feishu.cn/document/ukTMukTMukTM/ugDOyYjL4gjM24CO4IjN
	// 更多介绍请看：README.zh.md->如何构建配置（Config）
	appConf = lark.NewInternalAppConfig(lark.DomainFeiShu, 
		lark.SetAppCredentials("AppID", "AppSecret"), // 必需
		lark.SetAppEventKey("VerificationToken", "EncryptKey"), // 非必需，订阅事件、消息卡片时必需
		lark.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken")) // 非必需，使用服务台API时必需)
}

func test(ctx context.Context) {
	imService := im.NewService(appConf)
	
	larkCtx := lark.WrapContext(ctx)
	// 发送消息
	reqCall := imService.Messages.Create(larkCtx, &im.MessageCreateReqBody{
		ReceiveId: "ou_a11d2bcc7d852afbcaf37e5b3ad01f7e",
		Content:   `{"text":"test content"}`,
		MsgType:   "text",
	})
	reqCall.SetReceiveIdType("open_id")
	
	message, err := reqCall.Do()
	// 打印 request_id 方便 oncall 时排查问题
	fmt.Println(larkCtx.GetRequestID())
	fmt.Println(larkCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(message))
}

```

#### 使用`企业自建应用`访问 [发送文本消息API](https://open.feishu.cn/document/ukTMukTMukTM/uUjNz4SN2MjL1YzM) 示例

- 不在 [service](./service) 下的业务 API，没有直接可以使用的SDK，可以使用原生模式

```go
package demo

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
)

// 更多介绍请看：README.zh.md->如何构建配置（Config）
var appConf *lark.AppConfig

func test(ctx context.Context) {
	// 发送消息的内容
	body := map[string]interface{}{
		"open_id":  "ou_a11d2bcc7d852afbcaf37e5b3ad01f7e",
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": "test send message",
		},
	}
	// 请求的上下文
	larkCtx := lark.WrapContext(ctx)
	// 请求发送消息的结果
	ret := make(map[string]interface{})
	// 构建请求
	req := lark.NewAPIRequest("/open-apis/message/v4/send", "POST", lark.AccessTokenTypeTenant, body, &ret)
	
	// 发送请求
	err := lark.SendAPIRequest(larkCtx, appConf, req)
	// 打印请求的RequestID
	fmt.Println(larkCtx.GetRequestID())
	// 打印请求的响应状态吗
	fmt.Println(larkCtx.GetHTTPStatusCode())
	// 请求的error处理
	if err != nil {
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	// 打印请求的结果
	fmt.Println(lark.Prettify(ret))
}

```

### 订阅服务端事件

- **必看** [订阅事件概述](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM) ，了解订阅事件的过程及注意事项
- SDK 已经封装了 WebHook 地址的 CHALLENGE 校验
- 更多使用示例，请看：[sample/event](sample/event)（含：结合 gin 的使用）

#### 使用`企业自建应用`订阅 [员工变更事件](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/user/events/updated) 示例

- 在 [service](./service) 下的业务 Event，都是可以直接使用 SDK

```go
package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"net/http"
)

// 更多介绍请看：README.zh.md->如何构建配置（Config）
var appConf *lark.AppConfig

func main() {
	
	// 设置用户数据变更事件处理者
	contact.SetUserUpdatedEventHandler(appConf, func(ctx *lark.Context, event *contact.UserUpdatedEvent) error {
		// 打印请求的Request ID，方便 oncall 排查问题
		fmt.Println(ctx.GetRequestID())
		// 打印事件
		fmt.Println(lark.Prettify(event))
		return nil
	})

	// 设置 "开发者后台" -> "事件订阅" 请求网址 URL：https://domain/webhook/event
	lark.WebHook.EventWebServeRouter("/webhook/event", appConf)
	// startup event http server, port: 8089
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}

```

#### 使用`企业自建应用` 订阅 [首次启用应用事件](https://open.feishu.cn/document/ukTMukTMukTM/uQTNxYjL0UTM24CN1EjN) 示例

- 有些不在 [service](./service) 下的业务 Event，没有直接可以使用的 SDK，可以使用原生模式

```go
package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"net/http"
)

// 更多介绍请看：README.zh.md->如何构建配置（Config）
var appConf *lark.AppConfig

func main() {
    // 设置 首次启用应用事件（app_open） 的处理
	lark.WebHook.SetEventHandler(appConf, "app_open", func(ctx *lark.Context, e map[string]interface{}) error {
		// 打印请求的Request ID
		fmt.Println(ctx.GetRequestID())
		// 打印事件
		fmt.Println(lark.Prettify(e))
		return nil
	})

	// 设置 "开发者后台" -> "事件订阅" 请求网址 URL：https://domain/webhook/event
	lark.WebHook.EventWebServeRouter("/webhook/event", appConf)
	// startup event http server, port: 8089
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

### 处理消息卡片回调

- **必看** [消息卡片开发流程](https://open.feishu.cn/document/ukTMukTMukTM/uAzMxEjLwMTMx4CMzETM) ，了解订阅事件的过程及注意事项
- SDK 已经封装了 WebHook 地址的 CHALLENGE 校验
- 更多使用示例，请看：[sample/card](sample/card) 

#### 使用`企业自建应用`处理消息卡片回调示例

```go
package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"net/http"
)

// 更多介绍请看：README.zh.md->如何构建配置（Config）
var appConf *lark.AppConfig

func main() {

	// 设置消息卡片的处理者
	lark.WebHook.SetCardActionHandler(appConf, func(ctx *lark.Context, cardAction *lark.CardAction) (interface{}, error) {
		// 打印消息卡片
		fmt.Println(lark.Prettify(cardAction))
		return nil, nil // 返回值：可以为 nil、新的消息卡片的 JSON 字符串
	})

	// 设置 "开发者后台" -> "应用功能" -> "机器人" 消息卡片请求网址：https://domain/webhook/card
	lark.WebHook.CardWebServeRouter("/webhook/card", appConf)
	// startup event http server, port: 8089
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

## 如何构建配置（Config）

---

- 应用信息的配置
- 日志接口（Logger）的实现，用于输出 SDK 处理过程中产生的日志，便于排查问题。
    - 可以使用业务系统的日志实现，请看示例代码：[sample/config.go](sample/config.go) 的 Logrus
- 存储接口（Store）的实现，用于保存访问凭证（app/tenant_access_token）、临时凭证(app_ticket）
    - 推荐使用Redis实现，请看示例代码：[sample/config.go](sample/config.go) 的 RedisStore
        - 减少获取 访问凭证 的次数，防止调用访问凭证 接口被限频。
        - 应用商品应用，接受开放平台下发的 app_ticket，会保存到存储中，所以存储接口（Store）的实现的实现需要支持分布式存储。

```go

// domain：URL域名地址，值范围：lark.DomainFeiShu / lark.DomainLarkSuite / 其他URL域名地址(string)
// 防止应用信息泄漏，配置环境变量中，变量（4个）说明：
// APP_ID："开发者后台" -> "凭证与基础信息" -> 应用凭证 App ID
// APP_SECRET："开发者后台" -> "凭证与基础信息" -> 应用凭证 App Secret
// VERIFICATION_TOKEN："开发者后台" -> "事件订阅" -> 事件订阅 Verification Token
// ENCRYPT_KEY："开发者后台" -> "事件订阅" -> 事件订阅 Encrypt Key
// HELP_DESK_ID: 服务台设置中心 -> ID
// HELP_DESK_TOKEN: 服务台设置中心 -> 令牌

// 企业自建应用的配置，通过环境变量获取应用配置
appConf = lark.NewInternalAppConfigByEnv(domain lark.Domain)

// 应用商店应用的配置，通过环境变量获取应用配置
appConf = lark.NewISVAppConfigByEnv(domain lark.Domain)

// 参数说明：
// domain：URL域名地址，值范围：lark.DomainFeiShu / lark.DomainLarkSuite / 其他URL域名地址(string)
// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）
// HelpDeskID、HelpDeskToken：服务台设置中心 -> ID、令牌

// 企业自建应用的配置
appConf = lark.NewInternalAppConfig(domain lark.Domain,
    lark.SetAppCredentials("AppID", "AppSecret"),               // 必需
    lark.SetAppEventKey("VerificationToken", "EncryptKey"),     // 非必需，事件订阅时必需
    lark.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken"), // 非必需，访问服务台 API 时必需
)
// 应用商店应用的配置
appConf = lark.NewISVAppConfig(domain lark.Domain,
    lark.SetAppCredentials("AppID", "AppSecret"),               // 必需
    lark.SetAppEventKey("VerificationToken", "EncryptKey"),     // 非必需，事件订阅时必需
    lark.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken"), // 非必需，访问服务台 API 时必需
)

// 设置日志的实现（lark.Logger接口），默认：控制台输出
appConf.SetLogger(log lark.Logger)
// 设置日志的级别，默认：ERROR 级别
appConf.SetLogLevel(lark.LogLevelDebug) // 日志级别设置成 "DEBUG"，可以看到更多日志，方便排查问题
// 设置存储的实现（lark.Store接口），默认：本地内存（sync.Map）
appConf.SetStore(store lark.Store)

```

## 如何构建请求（APIRequest）

---

- 使用原生模式访问 API，这时需要构建请求
- 更多示例，请看：[sample/api/api.go](sample/api/api.go)（含：文件的上传与下载）

```go
import (
    "github.com/larksuite/oapi-sdk-go"
)

// 参数说明：
// httpPath：API 路径
    // 例如：https://domain/open-apis/contact/v3/users/:user_id
    // 支持：域名之后的路径，则 httpPath："/open-apis/contact/v3/users/:user_id"（推荐）
    // 支持：全路径，则 httpPath："https://domain/open-apis/contact/v3/users/:user_id"
    // 支持： /open-apis/ 之后的路径，则 httpPath："contact/v3/users/:user_id"
// httpMethod: GET/POST/PUT/BATCH/DELETE
// accessTokenType：API使用哪种访问凭证，取值范围：lark.AccessTokenTypeApp/lark.AccessTokenTypeTenant/lark.AccessTokenTypeUser，例如：lark.AccessTokenTypeTenant
// input：请求体（可能是 lark.NewFormData()（例如：文件上传））,如果不需要请求体（例如一些GET请求），则传：nil
// output：响应体（output := response["data"])
// opts：扩展函数，一些不常用的参数封装，如下：
    // lark.SetPathParams(map[string]interface{}{"user_id": 4})：设置URL Path参数（有:前缀）值，当 httpPath="contact/v3/users/:user_id" 时，请求的URL="https://{domain}/open-apis/contact/v3/users/4"
    // lark.SetQueryParams(map[string]interface{}{"age":4, "types":[1,2]})：设置 URL query，会在url追加?age=4&types=1&types=2
    // lark.setResponseStream()，设置响应的是否是流，例如下载文件，这时：output的类型需要实现 io.Writer 接口
    // lark.SetNotDataField()，设置响应的是否 没有`data`字段，业务接口都是有`data`字段，所以不需要设置
    // lark.SetTenantKey("TenantKey")，以`应用商店应用`身份，表示使用`tenant_access_token`访问API，需要设置
    // lark.SetUserAccessToken("UserAccessToken")，表示使用`user_access_token`访问API，需要设置
    // lark.NeedHelpDeskAuth()，表示是服务台API，需要设置 config.AppSettings 的 help desk 信息
req := lark.NewAPIRequest(httpPath, httpMethod string, accessTokenType lark.AccessTokenType,
    input interface{}, output interface{}, opts ...lark.APIRequestOpt)
```

## 如何构建请求上下文（lark.Context）及常用方法

---

```go
import(
    "context"
    "github.com/larksuite/oapi-sdk-go"
)

// 参数说明：
// c：context.Context
// 返回值说明：
// ctx: 实现了 Go 的 context.Context，保存请求处理过程中的一些变量
larkCtx := lark.WrapContext(c context.Context)

// 获取请求的 Request ID，便于排查问题
requestId := larkCtx.GetRequestID()

// 获取请求的响应状态码
httpStatusCode := larkCtx.GetHTTPStatusCode()

// 获取请求的响应 header
header := larkCtx.GetHeader()

```

## 如何发送请求

---

- 由于 SDK 已经封装了 app_access_token、tenant_access_token 的获取，所以在调业务API的时候，不需要去获取 app_access_token、tenant_access_token。
-  如果业务接口需要使用 user_access_token，需要进行设置（lark.SetUserAccessToken("UserAccessToken")），具体请看 README.zh.md -> 如何构建请求（Request）
- 更多使用示例，请看：[sample/api/api.go](sample/api/api.go)

```go
import(
    "github.com/larksuite/oapi-sdk-go"
)

// 参数说明：
// ctx：请求的上下文（）
// conf：配置（Config）
// req：请求（APIRequest）
// 返回值说明：
// err：发送请求，出现的错误以及响应的错误码（response.body["code"]）不等于0
err := lark.SendAPIRequest(ctx *lark.Context, conf lark.Config, req *lark.APIRequest)

```

## 下载文件工具

---

- 通过网络请求下载文件
- 更多使用示例，请看：[sample/tools/file_download.go](sample/tools/file_download.go)

```go
import(
    "context"
    "github.com/larksuite/oapi-sdk-go"
)

// 获取文件内容
// 参数说明：
// ctx：context.Context
// url：文件的HTTP地址
// 返回值说明：
// bytes：文件内容的二进制数组
// err：错误
bytes, err := lark.DownloadFile(ctx context.Context, url string)

// 获取文件内容流，读取完文件内容后，需要关闭流
// 参数说明：
// ctx：context.Context
// url：文件的HTTP地址
// 返回值说明：
// readCloser：文件内容的二进制读取流
// err：错误
readCloser, err := lark.DownloadFileToStream(ctx context.Context, url string)

```

## License

---

- MIT



