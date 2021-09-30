# 飞书开放接口SDK

## 概述

- 飞书开放平台，便于企业应用与飞书集成，让协同与管理更加高效
- 飞书开发接口SDK，便捷 [调用服务端API](https://open.feishu.cn/document/ukTMukTMukTM/uITNz4iM1MjLyUzM)
  与 [订阅服务端事件](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM)。

## 问题反馈

如有任何SDK使用相关问题，请提交 [Github Issues](https://github.com/larksuite/oapi-sdk-go/issues), 我们会在收到 Issues 的第一时间处理，并尽快给您答复。

## 运行环境

- Golang 1.5及以上

## 安装方法

```shell
go get github.com/larksuite/v2@v2.0.0-rc1
```

## SDK 包引入与使用规则

- lark 包，引入的路径："github.com/larksuite/oapi-sdk-go/v2"

- 业务 包，引入的路径："github.com/larksuite/oapi-sdk-go/v2/service/业务/版本"

    - 例如：im 包，引入的路径："github.com/larksuite/oapi-sdk-go/v2/service/im/v1"

- SDK 包如何使用，下面有代码示例可以参考

## 术语解释

- [开发文档](https://open.feishu.cn/document/) ：开放平台的开放接口的参考，**开发者必看，可以使用搜索功能，高效的查询文档**
- [开发者后台](https://open.feishu.cn/app/) ：开发者开发应用的管理后台
- [企业自建应用](https://open.feishu.cn/document/home/introduction-to-custom-app-development/self-built-application-development-process)
  ：应用仅仅可在本企业内发布使用
- [应用商店应用](https://open.feishu.cn/document/uMzNwEjLzcDMx4yM3ATM/uYzNwEjL2cDMx4iN3ATM?lang=zh-CN)
  ：应用会在 [应用目录](https://app.feishu.cn/) 展示，各个企业可以选择安装使用

![App type](doc/app_type.zh.png)

## 如何调用服务端API

- **必看** [如何调用服务端 API](https://open.feishu.cn/document/ukTMukTMukTM/uYTM5UjL2ETO14iNxkTN/guide-to-use-server-api)
  ，了解调用服务端API的过程及注意事项。
    - 由于 SDK 已经封装了 app_access_token、tenant_access_token 的获取，所以在调业务 API 的时候，不需要去获取
      app_access_token、tenant_access_token。如果业务接口需要使用 user_access_token，需要进行设置 lark.WithUserAccessToken("userAccessToken")，具体请看：[如何发送请求](#如何发送请求)

### 使用`企业自建应用`访问 [发送消息 API](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/create) 示例

- 在 [v2/service](./v2/service) 下的业务 API，都是可以直接使用业务 SDK
- 更多示例，请看：[v2/sample/api/im.go](./v2/sample/api/im.go)（含：文件的上传与下载）

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/service/im/v1"
	"os"
)

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")

	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret)

	ctx := context.Background()
	messageCreate(ctx, larkApp)
}

// 发送消息
func messageCreate(ctx context.Context, larkApp *lark.App) {
	messageText := &lark.MessageText{Text: "Tom test content"}
	content, err := messageText.JSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	messageCreateResp, err := im.New(larkApp).Messages.Create(ctx, &im.MessageCreateReq{
		ReceiveIdType: lark.StringPtr("user_id"),
		Body: &im.MessageCreateReqBody{
			ReceiveId: lark.StringPtr("77bbc392"),
			MsgType:   lark.StringPtr("text"),
			Content:   lark.StringPtr(content),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", messageCreateResp.RequestId())
	if messageCreateResp.Code != 0 {
		fmt.Println(messageCreateResp.CodeError)
		return
	}
	fmt.Println(lark.Prettify(messageCreateResp.Data))
}
```

### 使用`企业自建应用`访问 [发送文本消息 API](https://open.feishu.cn/document/ukTMukTMukTM/uUjNz4SN2MjL1YzM) 示例

- 有些老版接口，没有直接可以使用的业务 SDK，可以使用`原生`模式，具体请看：[如何发送请求](#如何发送请求)

- 更多示例，请看：[v2/sample/api/api.go](./v2/sample/api/api.go)（含：文件的上传与下载）

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"net/http"
	"os"
)

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")

	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret)

	ctx := context.Background()
	sendMessage(ctx, larkApp)
}

// 发送消息
func sendMessage(ctx context.Context, larkApp *lark.App) {
	resp, err := larkApp.SendRequest(ctx, http.MethodGet, "/open-apis/message/v4/send", 
		lark.AccessTokenTypeTenant, map[string]interface{}{
		"user_id":  "77bbc392",
		"msg_type": "text",
		"content":  &lark.MessageText{Text: "test"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", resp.RequestId())
	fmt.Println(resp)
	fmt.Println()
}

```

## 如何订阅服务端事件

- **必看** [订阅事件概述](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM) ，了解订阅事件的过程及注意事项

### 使用`企业自建应用`订阅 [接收消息事件](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/events/receive) 示例

- 在 [v2/service](./v2/service) 下的业务 Event，都是可以直接使用业务 SDK
  
- 更多使用示例，请看：[v2/sample/event/im.go](./v2/sample/event/im.go)

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/service/im/v1"
	"net/http"
	"os"
)

func main() {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY")

	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithAppEventVerify(verificationToken, encryptKey))
	
	// @应用机器人的消息处理
	im.New(larkApp).Messages.ReceiveEventHandler(func(ctx context.Context, req *lark.RawRequest, event *im.MessageReceiveEvent) error {
		fmt.Println(req)
		fmt.Println(lark.Prettify(event))
		return nil
	})

	// http server handle func
	http.HandleFunc("/webhook/event", func(writer http.ResponseWriter, request *http.Request) {
		// 如果开发者使用是其他 Web 框架，需要将 Web 框架的 Request 装成 lark.RawRequest
		// 经过 larkApp.Webhook.EventCommandHandle(...) 的处理，返回 lark.RawResponse
		// 再将 lark.RawResponse 转成 Web 框架的Response 
		rawRequest, err := lark.NewRawRequest(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		rawResp := larkApp.Webhook.EventCommandHandle(context.Background(), rawRequest)
		rawResp.Write(writer)
	})
	// 设置 "开发者后台" -> "事件订阅" 请求网址 URL：https://domain/webhook/event
	// startup event http server, port: 8089
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

### 使用`企业自建应用`订阅 [接收消息事件](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/events/receive) 示例

- 有些老的事件，没有直接可以使用的业务 SDK，可以使用`原生`模式，具体请看：[如何订阅事件](#如何订阅事件)
  
- 更多使用示例，请看：[v2/sample/event/event.go](./v2/sample/event/event.go)

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"net/http"
	"os"
)

func main() {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY")

	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithAppEventVerify(verificationToken, encryptKey))
	
	// @应用机器人的消息处理
	// "im.message.receive_v1"：事件类型
	larkApp.Webhook.EventHandleFunc("im.message.receive_v1", func(ctx context.Context, req *lark.RawRequest) error {
		fmt.Println(req.RequestId())
		fmt.Println(req)
		return nil
	})

	// http server handle func
	http.HandleFunc("/webhook/event", func(writer http.ResponseWriter, request *http.Request) {
		// 如果开发者使用是其他 Web 框架，需要将 Web 框架的 Request 装成 lark.RawRequest
		// 经过 larkApp.Webhook.EventCommandHandle(...) 的处理，返回 lark.RawResponse
		// 再将 lark.RawResponse 转成 Web 框架的Response
		rawRequest, err := lark.NewRawRequest(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		rawResp := larkApp.Webhook.EventCommandHandle(context.Background(), rawRequest)
		rawResp.Write(writer)
	})
	// 设置 "开发者后台" -> "事件订阅" 请求网址 URL：https://domain/webhook/event
	// startup event http server, port: 8089
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

## 如何处理消息卡片 Action

- **必看** [消息卡片开发流程](https://open.feishu.cn/document/ukTMukTMukTM/uAzMxEjLwMTMx4CMzETM) ，了解订阅事件的过程及注意事项
- 更多使用示例，请看：[v2/sample/card/card.go](./v2/sample/card/card.go)

#### 使用`企业自建应用`处理消息卡片回调示例

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"net/http"
	"os"
)

func main() {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY")
	
	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithAppEventVerify(verificationToken, encryptKey))

	// card action handler
	// return new card
	larkApp.Webhook.CardActionHandler(func(ctx context.Context, request *lark.RawRequest,
		action *lark.CardAction) (interface{}, error) {
		fmt.Println(request)
		fmt.Println(lark.Prettify(action))

		card := &lark.MessageCard{
			Config: &lark.MessageCardConfig{WideScreenMode: lark.BoolPtr(true)},
			Elements: []lark.MessageCardElement{&lark.MessageCardMarkdown{
				Content: "**test**",
			}},
		}
		return card, nil
	})
	// http server handle func
	http.HandleFunc("/webhook/card", func(writer http.ResponseWriter, request *http.Request) {
		// 如果开发者使用是其他 Web 框架，需要将 Web 框架的 Request 装成 *lark.RawRequest
		// 经过 larkApp.Webhook.CardActionHandle(...) 的处理，返回 *lark.RawResponse
		// 再将 lark.RawResponse 转成 Web 框架的Response
		rawRequest, err := lark.NewRawRequest(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		larkApp.Webhook.CardActionHandle(context.Background(), rawRequest).Write(writer)
	})
	// 设置 "开发者后台" -> "应用功能" -> "机器人" 消息卡片请求网址：https://domain/webhook/card
	// startup event http server, port: 8089
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

## 使用`应用商店应用`调用 服务端API 示例

- [如何获取 app_access_token](https://open.feishu.cn/document/ukTMukTMukTM/uEjNz4SM2MjLxYzM) （应用商店应用）
    - 与企业自建应用相比，应用商店应用的获取 app_access_token 的流程复杂一些。
        - 需要开放平台下发的 app_ticket，通过订阅事件接收。SDK 已经封装了 app_ticket 事件的处理，只需要启动事件订阅服务。
        - 使用SDK调用服务端 API 时，如果当前还没有收到开发平台下发的 app_ticket ，会报错且向开放平台申请下发 app_ticket ，可以尽快的收到开发平台下发的 app_ticket，保证再次调用服务端 API 的正常。
        - 使用SDK调用服务端 API 时，需要使用 tenant_access_token 访问凭证时，需要 tenant_key ，来表示当前是哪个租户使用这个应用调用服务端 API。
            - tenant_key，租户安装启用了这个应用，开放平台发送的服务端事件，事件内容中都含有 tenant_key。
- 示例代码：[v2/sample/api/marketplace_app.go](./v2/sample/api/marketplace_app.go)

## 如何构建应用

```go

import (
    "github.com/larksuite/oapi-sdk-go/v2"
)

// 防止应用信息泄漏，配置环境变量中，变量说明：
// APP_ID："开发者后台" -> "凭证与基础信息" -> 应用凭证 App ID
// APP_SECRET："开发者后台" -> "凭证与基础信息" -> 应用凭证 App Secret
// VERIFICATION_TOKEN："开发者后台" -> "事件订阅" -> 事件订阅 Verification Token （事件订阅、处理消息卡片Action 必需）
// ENCRYPT_KEY："开发者后台" -> "事件订阅" -> 事件订阅 Encrypt Key（事件订阅 必需）
// HELP_DESK_ID: 服务台设置中心 -> ID
// HELP_DESK_TOKEN: 服务台设置中心 -> 令牌

appID, appSecret, verificationToken, encryptKey, helpDeskID, helpDeskToken := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY"), os.Getenv("HELP_DESK_ID"), os.Getenv("HELP_DESK_TOKEN")

// 企业自建应用的配置
larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
    lark.WithAppEventVerify(verificationToken, encryptKey), // 非必需，事件订阅、处理消息卡片Action时必需
    lark.WithAppHelpdeskCredential(helpDeskID, helpDeskToken), // 非必需，访问服务台API时必需
)

// 应用商店应用的配置
larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret, 
    lark.WithAppType(lark.AppTypeMarketplace), // 标识应用类型为：应用商店应用
    lark.WithAppEventVerify(verificationToken, encryptKey), // 非必需，事件订阅、处理消息卡片Action时必需
    lark.WithAppHelpdeskCredential(helpDeskID, helpDeskToken), // 非必需，访问服务台API时必需
)


// 配置日志接口的实现
// 例如：使用 logrus 实现，请看示例代码：v2/sample/logrus.go
// 例如：日志（lark.NewDefaultLogger()：日志控制台输出），日志级别（lark.LogLevelDebug：debug级别，可以打印更好的日志，利于排查问题）
larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
	lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug), // 非必需
)
// 更多示例：v2/sample/api/marketplace_app.go 的 "sample.Logrus{}"
larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
    lark.WithLogger(sample.Logrus{}, lark.LogLevelDebug),
)


// 配置存储接口，用于存放：app_access_token、tenant_access_token、app_ticket
// 默认是 sync.map 内存实现的
// 例如：使用 Redis 实现，请看示例代码：v2/sample/redis_store.go
// 对于应用商品应用，接收开放平台下发的 app_ticket，会保存到存储中，所以存储接口的实现的实现需要支持分布式存储
// 更多示例：v2/sample/api/marketplace_app.go 的 "sample.NewRedisStore()"
larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
	lark.WithStore(sample.NewRedisStore()) // use redis store
)

```

## 如何发送请求

- 有些老版接口，没有直接可以使用的业务 SDK，可以使用原生模式
- 更多示例，请看：[v2/sample/api/api.go](./v2/sample/api/api.go)（含：文件的上传与下载）

```go

import (
    "context"
    "net/http"
    "github.com/larksuite/oapi-sdk-go/v2"
)

app := lark.NewApp(lark.DomainFeiShu, appID, appSecret)


// 参数说明：
// ctx: context.Context

// httpMethod: HTTP method（http.MethodGet/http.MethodPost/http.MethodPut/http.MethodPatch/http.MethodDelete）

// httpPath：API路径
// 支持：域名之后的路径，则 httpPath："/open-apis/contact/v3/users/:user_id"（推荐）
// 支持：全路径，则 httpPath："https://domain/open-apis/contact/v3/users/:user_id"

// accessTokenType：API使用哪种访问凭证（lark.AccessTokenTypeApp/lark.AccessTokenTypeTenant/lark.AccessTokenTypeUser）

// input：请求体（可能是 lark.NewFormdata()（例如：文件上传））, 如果不需要请求体（例如：GET请求），则传：nil

// options：扩展函数，如下： 
// lark.WithTenantKey("tenantKey")，以`应用商店应用`身份，表示使用`tenant_access_token`访问API，需要设置
// lark.WithUserAccessToken("userAccessToken")，表示使用`user_access_token`访问API，需要设置
// lark.WithNeedHelpDeskAuth()，表示是服务台API，需要HelpDesk的Auth验证，需要设置 lark app 的 HelpDesk 信息
// lark.WithHTTPHeader(header http.Header)，设置 HTTP header
func (app *App) SendRequest(ctx context.Context, httpMethod string, httpPath string,
    accessTokenType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {}


// 发送请求的响应（lark.RawResponse）
type RawResponse struct {
    StatusCode int         `json:"-"`
    Header     http.Header `json:"-"`
    RawBody    []byte      `json:"-"`
}

// 获取请求的ID，反馈问题的时候，提供RequestId（HTTP.header["X-Tt-Logid"]），排查问题更方便
func (resp RawResponse) RequestId() string {}

// 响应的Body，反序列化一个实例上
func (resp RawResponse) JSONUnmarshalBody(val interface{}) error {}

```

## 如何订阅事件

- 有些老版接口，没有直接可以使用的业务 SDK，可以使用原生模式
- 更多示例，请看：[v2/sample/event/event.go](./v2/sample/event/event.go)

```go

import (
    "context"
    "net/http"
    "github.com/larksuite/oapi-sdk-go/v2"
)

app := lark.NewApp(lark.DomainFeiShu, appID, appSecret, lark.WithAppEventVerify(verificationToken, encryptKey))

// app.Webhook.EventHandleFunc
// 参数说明：
// eventType: 事件类型

// handler: 事件处理函数
// ctx：context.Context
// req：事件的回调请求
func (wh *webhook) EventHandleFunc(eventType string, handler func(ctx context.Context, req *RawRequest) error) {}


// 事件的回调请求（lark.RawRequest）
type RawRequest struct {
    Header  http.Header
    RawBody []byte
}

// 获取请求的ID，反馈问题的时候，提供RequestId，排查问题更方便
func (req RawRequest) RequestId() string {}

// 请求的Body，反序列化一个实例上
func (req RawRequest) JSONUnmarshalBody(val interface{}) error {}

```

## 消息内容 Model

- 文档：[发送消息 content 说明](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/im-v1/message/create_json)

- 消息内容 Model 代码：[v2/message_model.go](./v2/message_model.go)

- 消息内容 Model 使用示例：[v2/sample/api/message_model.go](./v2/sample/api/message_model.go)

|消息类型| Model |
|----|----|
|文本 text|lark.MessageText|
|富文本 post|lark.MessagePost|
|图片 image|lark.MessageImage|
|消息卡片 interactive|lark.MessageCard|
|分享个人名片 share_user|lark.MessageShareUser|
|分享群组名片 share_chat|lark.MessageShareChat|
|语音 audio|lark.MessageAudio|
|视频 media|lark.MessageVideo|
|文件 file|lark.MessageFile|
|表情包 sticker|lark.MessageFile|

```go

// 消息内容Model都有JSON方法，返回JSON字符串
func (m *Message***) JSON() (string, error) {}

```

## 基本类型与指针类型的转换

### 基本类型转指针类型
|方法名| 描述 |
|----|----|
|lark.StringPtr(v string)|string 转 *string|
|lark.BoolPtr(v bool)|bool 转 *bool|
|lark.IntPtr(v int)|int 转 *int|
|lark.Int8Ptr(v int8)|int8 转 *int8|
|lark.Int16Ptr(v int16)|int16 转 *int16|
|lark.Int32Ptr(v int32)|int32 转 *int32|
|lark.Int64Ptr(v int64)|int64 转 *int64|
|lark.Float32Ptr(v float32)|float32 转 *float32|
|lark.Float64Ptr(v float64)|float64 转 *float64|
|lark.TimePtr(v time.Time)|time.Time 转 *time.Time|

### 指针类型转基本类型
|方法名| 描述 |
|----|----|
|lark.StringValue(v *string)|*string 转 string|
|lark.BoolValue(v *bool)|*bool 转 bool|
|lark.IntValue(v *int)|*int 转 int|
|lark.Int8Value(v *int8)|*int8 转 int8|
|lark.Int16Value(v *int16)|*int16 转 int16|
|lark.Int32Value(v *int32)|*int32 转 int32|
|lark.Int32Value(v *int64)|*int64 转 int64|
|lark.Float32Value(v *float32)|*float32 转 float32|
|lark.Float64Value(v *float64)|*float64 转 float64|
|lark.TimeValue(v *time.Time)|*time.Time 转 time.Time|


## 下载文件工具

- 通过网络请求下载文件
- 更多使用示例，请看：[v2/sample/utils/file_download.go](./v2/sample/utils/file_download.go)

## License

---

- MIT



