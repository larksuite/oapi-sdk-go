[**README of Larksuite(Overseas)**](README.md) | 飞书

# 飞书开放接口SDK

## 概述

---

- 飞书开放平台，便于企业应用与飞书集成，让协同与管理更加高效，[概述](https://open.feishu.cn/document/uQjL04CN/ucDOz4yN4MjL3gzM)

- 飞书开发接口SDK，便捷调用服务端API与订阅服务端事件，例如：消息&群组、通讯录、日历、视频会议、云文档、 OKR等具体可以访问 [飞书开放平台文档](https://open.feishu.cn/document/) 看看【服务端
  API】。

## 运行环境

---

- Golang 1.5及以上

## 安装方法

---

```shell
go get -u github.com/larksuite/oapi-sdk-go
```

## 术语解释
- 飞书（FeiShu）：Lark在中国的称呼，主要为国内的企业提供服务，拥有独立的[域名地址](https://www.feishu.cn)。
- LarkSuite：Lark在海外的称呼，主要为海外的企业提供服务，拥有独立的[域名地址](https://www.larksuite.com/) 。
- 开发文档：开放平台的开放接口的参考，**开发者必看，可以使用搜索高效的查询文档**。[更多介绍说明](https://open.feishu.cn/document/) 。
- 开发者后台：开发者开发应用的管理后台，[更多介绍说明](https://open.feishu.cn/app/) 。
- 企业自建应用：应用仅仅可在本企业内安装使用，[更多介绍说明](https://open.feishu.cn/document/uQjL04CN/ukzM04SOzQjL5MDN) 。
- 应用商店应用：应用会在 [应用目录](https://app.feishu.cn/?lang=zh-CN) 展示，各个企业可以选择安装，[更多介绍说明](https://open.feishu.cn/document/uQjL04CN/ugTO5UjL4kTO14CO5kTN) 。
  
![App type](doc/app_type.zh.png)

## 快速使用

---

### 调用服务端API

- **必看** [如何调用服务端API](https://open.feishu.cn/document/ukTMukTMukTM/uYTM5UjL2ETO14iNxkTN/guide-to-use-server-api)
  ，了解调用服务端API的过程及注意事项。

#### 使用`企业自建应用`访问 发送文本消息API 示例

- 有些老版接口，没有直接可以使用的SDK，可以使用`原生`模式。

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/tools"
)

var conf *config.Config

func init() {
	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）。
	appSetting := config.NewInternalAppSettings("AppID", "AppSecret", "VerificationToken", "EncryptKey")

	// 当前访问的是飞书，使用默认存储、默认日志（Debug级别），更多可选配置，请看：README.zh.md->高级使用->如何构建整体配置（Config）。
	conf = config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelInfo)
}

func main() {
	// 发送消息的内容
	body := map[string]interface{}{
		"open_id":  "user open id",
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": "test send message",
		},
	}
	// 请求发送消息的结果
	ret := make(map[string]interface{})
	// 构建请求
	req := request.NewRequestWithNative("message/v4/send", "POST", request.AccessTokenTypeTenant, body, &ret)
	// 请求的上下文
	coreCtx := core.WrapContext(context.Background())
	// 发送请求
	err := api.Send(coreCtx, conf, req)
	// 打印请求的RequestID
	fmt.Println(coreCtx.GetRequestID())
	// 打印请求的响应状态吗
	fmt.Println(coreCtx.GetHTTPStatusCode())
	// 请求的error处理
	if err != nil {
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		fmt.Println(tools.Prettify(err))
		return
	}
	// 打印请求的结果
	fmt.Println(tools.Prettify(ret))
}
```

#### 使用`企业自建应用`访问 修改用户部分信息API 示例

- 该接口是新的接口，可以直接使用SDK。

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
)

var conf *config.Config

func init() {
	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）。
	appSetting := config.NewInternalAppSettings("AppID", "AppSecret", "VerificationToken", "EncryptKey")

	// 当前访问的是飞书，使用默认存储、默认日志（Debug级别），更多可选配置，请看：README.zh.md->高级使用->如何构建整体配置（Config）。
	conf = config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelInfo)
}

func main() {
	service := contact.NewService(conf)
	coreCtx := core.WrapContext(context.Background())
	body := &contact.User{}
	body.Name = "rename"
	// 由于这是一个PATCH请求，需要明确更新哪些字段
	body.ForceSendFields = append(body.ForceSendFields, "Name")
	reqCall := service.Users.Patch(coreCtx, body)
	reqCall.SetUserId("user id")
	reqCall.SetUserIdType("user_id")
	// 发送请求
	result, err := reqCall.Do()
	// 打印请求的RequestID
	fmt.Println(coreCtx.GetRequestID())
	// 打印请求的响应状态吗
	fmt.Println(coreCtx.GetHTTPStatusCode())
	// 请求的error处理
	if err != nil {
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		fmt.Println(tools.Prettify(err))
		return
	}
	// 打印请求的结果
	fmt.Println(tools.Prettify(result))
}
```

#### [使用`应用商店应用`调用 服务端API 示例](doc/ISV.APP.README.zh.md)

### 订阅服务端事件

- **必看** [订阅事件概述](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM) ，了解订阅事件的过程及注意事项。
- 更多使用示例，请看[sample/event](sample/event)（含：结合gin的使用）

#### 使用`企业自建应用` 订阅 [首次启用应用事件](https://open.feishu.cn/document/ukTMukTMukTM/uQTNxYjL0UTM24CN1EjN) 示例

- 有些老的事件，没有直接可以使用的SDK，可以使用`原生`模式

```go
package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/event"
	eventhttpserver "github.com/larksuite/oapi-sdk-go/event/http/native"
	"net/http"
)

var conf *config.Config

func init() {
	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）。
	appSetting := config.NewInternalAppSettings("AppID", "AppSecret", "VerificationToken", "EncryptKey")

	// 当前访问的是飞书，使用默认存储、默认日志（Debug级别），更多可选配置，请看：README.zh.md->高级使用->如何构建整体配置（Config）。
	conf = config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelInfo)
}

func main() {
	// 设置首次启用应用事件callback
	event.SetTypeCallback(conf, "app_open", func(ctx *core.Context, e map[string]interface{}) error {
		// 打印请求的Request ID
		fmt.Println(ctx.GetRequestID())
		// 打印事件
		fmt.Println(tools.Prettify(e))
		return nil
	})

	// 启动httpServer，"开发者后台" -> "事件订阅" 请求网址 URL：https://domain/webhook/event
	eventhttpserver.Register("/webhook/event", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

#### 使用`企业自建应用`订阅 [用户数据变更事件](https://open.feishu.cn/document/ukTMukTMukTM/uITNxYjLyUTM24iM1EjN#70402aa) 示例

- 该接口是新的事件，可以直接使用SDK

```go
package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	eventhttpserver "github.com/larksuite/oapi-sdk-go/event/http/native"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"net/http"
)

var conf *config.Config

func init() {
	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）。
	appSetting := config.NewInternalAppSettings("AppID", "AppSecret", "VerificationToken", "EncryptKey")

	// 当前访问的是飞书，使用默认存储、默认日志（Debug级别），更多可选配置，请看：README.zh.md->高级使用->如何构建整体配置（Config）。
	conf = config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelInfo)
}

func main() {
	// 设置用户数据变更事件处理者
	contact.SetUserUpdatedEventHandler(conf, func(ctx *core.Context, event *contact.UserUpdatedEvent) error {
		// 打印请求的Request ID
		fmt.Println(ctx.GetRequestID())
		// 打印事件
		fmt.Println(tools.Prettify(event))
		return nil
	})

	// 启动httpServer，"开发者后台" -> "事件订阅" 请求网址 URL：https://domain/webhook/event
	// startup event http server, port: 8089
	eventhttpserver.Register("/webhook/event", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

### 处理消息卡片回调

- **必看** [消息卡片开发流程](https://open.feishu.cn/document/ukTMukTMukTM/uAzMxEjLwMTMx4CMzETM) ，了解订阅事件的过程及注意事项
- 更多使用示例，请看：[sample/card](sample/card) （含：结合gin的使用）

#### 使用`企业自建应用`处理消息卡片回调示例

```go
package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/card"
	cardhttpserver "github.com/larksuite/oapi-sdk-go/card/http/native"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"net/http"
)

var conf *config.Config

func init() {
	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）。
	appSetting := config.NewInternalAppSettings("AppID", "AppSecret", "VerificationToken", "EncryptKey")

	// 当前访问的是飞书，使用默认存储、默认日志（Debug级别），更多可选配置，请看：README.zh.md->高级使用->如何构建整体配置（Config）。
	conf = config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelInfo)
}

func main() {
	// 设置消息卡片的处理者
	// 返回值：可以为nil、新的消息卡片的Json字符串 
	card.SetHandler(conf, func(ctx *core.Context, c *model.Card) (interface{}, error) {
		// 打印消息卡片
		fmt.Println(tools.Prettify(c))
		return "{\"config\":{\"wide_screen_mode\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"div\",\"text\":{\"tag\":\"lark_md\",\"content\":\"[飞书golang](https://www.feishu.cn)整合即时沟通、日历、音视频会议、云文档、云盘、工作台等功能于一体，成就组织和个人，更高效、更愉悦。\"}}]}}", nil
	})
	// 设置 "开发者后台" -> "应用功能" -> "机器人" 消息卡片请求网址：https://domain/webhook/card
	// startup event http server, port: 8089
	cardhttpserver.Register("/webhook/card", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```    

## 高级使用

---

### 如何构建应用配置（AppSettings）

```go
import (
    "github.com/larksuite/oapi-sdk-go/core/config"
)

// 防止应用信息泄漏，配置环境变量中，变量（4个）说明：
// APP_ID："开发者后台" -> "凭证与基础信息" -> 应用凭证 App ID
// APP_SECRET："开发者后台" -> "凭证与基础信息" -> 应用凭证 App Secret
// VERIFICATION_TOKEN："开发者后台" -> "事件订阅" -> 事件订阅 Verification Token
// ENCRYPT_KEY："开发者后台" -> "事件订阅" -> 事件订阅 Encrypt Key
// 企业自建应用的配置，通过环境变量获取应用配置
appSettings := config.GetInternalAppSettingsByEnv()
// 应用商店应用的配置，通过环境变量获取应用配置
appSettings := config.GetISVAppSettingsByEnv()


// 参数说明：
// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）
// 企业自建应用的配置
appSettings := config.NewInternalAppSettings(appID, appSecret, verificationToken, encryptKey string)
// 应用商店应用的配置
appSettings := config.NewISVAppSettings(appID, appSecret, verificationToken, encryptKey string)

```

### 如何构建整体配置（Config）

- 访问 飞书、LarkSuite或者其他
- 应用的配置
- 日志接口（Logger）的实现，用于输出SDK处理过程中产生的日志，便于排查问题。
    - 可以使用业务系统的日志实现，请看示例代码：[sample/config/logrus.go](sample/config/logrus.go)
- 存储接口（Store）的实现，用于保存访问凭证（app/tenant_access_token）、临时凭证(app_ticket）
    - 推荐使用Redis实现，请看示例代码：[sample/config/redis_store.go](sample/config/redis_store.go)
        - 减少获取 访问凭证 的次数，防止调用访问凭证 接口被限频。
        - 应用商品应用，接受开放平台下发的app_ticket，会保存到存储中，所以存储接口（Store）的实现的实现需要支持分布式存储。

```go
import (
    "github.com/larksuite/oapi-sdk-go/core/config"
    "github.com/larksuite/oapi-sdk-go/core/constants"
    "github.com/larksuite/oapi-sdk-go/core/log"
    "github.com/larksuite/oapi-sdk-go/core/store"
)

// 方法一，推荐使用Redis实现存储接口（Store），减少访问获取AccessToken接口的次数
// 参数说明：
// domain：URL域名地址，值范围：constants.DomainFeiShu / constants.DomainLarkSuite / 其他URL域名地址
// appSettings：应用配置
// logger：[日志接口](core/log/log.go)
// loggerLevel：输出的日志级别 log.LevelInfo/LevelInfo/LevelWarn/LevelError
// store: [存储接口](core/store/store.go)，用来存储 app_ticket/access_token
conf := config.NewConfig(domain constants.Domain, appSettings *AppSettings, logger log.Logger, logLevel log.Level, store store.Store)

// 方法二，使用默认的存储接口（Store）的实现，适合轻量的使用（不合适：应用商店应用或调用服务端API次数频繁）
// 参数说明：
// domain：constants.DomainFeiShu / constants.DomainLarkSuite / 其他域名地址
// appSettings：应用配置
// logger：[日志接口](core/log/log.go)
// loggerLevel：输出的日志级别 log.LevelInfo/LevelInfo/LevelWarn/LevelError
conf := config.NewConfig(domain constants.Domain, appSettings *AppSettings, logger log.Logger, logLevel log.Level)

```

### 如何构建请求（Request）

- 有些老版接口，没有直接可以使用的SDK，可以使用原生模式，这时需要构建请求。
- 更多示例，请看：[sample/api/api.go](sample/api/api.go)（含：文件的上传与下载）

```go
import (
    "github.com/larksuite/oapi-sdk-go/api/core/request"
)

// 参数说明：
// httpPath：API路径（`open-apis/`之后的路径），例如：https://domain/open-apis/contact/v3/users/:user_id，则 httpPath："contact/v3/users/:user_id"
// httpMethod: GET/POST/PUT/BATCH/DELETE
// accessTokenType：API使用哪种访问凭证，取值范围：request.AccessTokenTypeApp/request.AccessTokenTypeTenant/request.AccessTokenTypeUser，例如：request.AccessTokenTypeTenant
// input：请求体（可能是request.NewFormData()（例如：文件上传））,如果不需要请求体（例如一些GET请求），则传：nil
// output：响应体（output := response["data"]) 
// optFns：扩展函数，一些不常用的参数封装，如下：
// request.SetPathParams(map[string]interface{}{"user_id": 4})：设置URL Path参数（有:前缀）值，当httpPath="contact/v3/users/:user_id"时，请求的URL="https://{domain}/open-apis/contact/v3/users/4"
// request.SetQueryParams(map[string]interface{}{"age":4,"types":[1,2]})：设置 URL qeury，会在url追加?age=4&types=1&types=2      
// request.setResponseStream()，设置响应的是否是流，例如下载文件，这时：output值是Buffer类型
// request.SetNotDataField(),设置响应的是否 没有`data`字段，业务接口都是有`data`字段，所以不需要设置
// request.SetTenantKey("TenantKey")，以`应用商店应用`身份，表示使用`tenant_access_token`访问API，需要设置
// request.SetUserAccessToken("UserAccessToken")，表示使用`user_access_token`访问API，需要设置
req := request.NewRequestWithNative(httpPath, httpMethod string, accessTokenType AccessTokenType, input interface{}, output interface{}, optFns ...OptFn)

```

### 如何构建请求上下文（core.Context）及常用方法

```go
import(
    "github.com/larksuite/oapi-sdk-go/core"
    "github.com/larksuite/oapi-sdk-go/core/config"
)

// 参数说明：
// c：context.Context
// 返回值说明：
// ctx: 实现了Golang的context.Context，保存请求中的一些变量
ctx := core.WrapContext(c context.Context)

// 获取请求的Request ID，便于排查问题
requestId := ctx.GetRequestID()

// 获取请求的响应状态码
httpStatusCode := ctx.GetHTTPStatusCode()

// 在事件订阅与消息卡片回调的处理者中，可以从core.Context中获取 Config
conf := config.ByCtx(ctx *core.Context)

```

### 如何发送请求

- 更多使用示例，请看：[sample/api/api.go](sample/api/api.go)

```go
import(
    "fmt"
    "context"
    "github.com/larksuite/oapi-sdk-go/api"
    "github.com/larksuite/oapi-sdk-go/api/core/request"
    "github.com/larksuite/oapi-sdk-go/api/core/response"
    "github.com/larksuite/oapi-sdk-go/core"
    "github.com/larksuite/oapi-sdk-go/core/test"
    "github.com/larksuite/oapi-sdk-go/core/tools"
)

// 参数说明：
// ctx：请求的上下文
// conf：整体的配置（Config）
// req：请求（Request）
// 返回值说明：
// err：发送请求，出现的错误以及响应的错误码（response.body["code"]）不等于0
err := api.Send(ctx *core.Context, conf *config.Config, req *request.Request)

```

### 下载文件工具

- 通过网络请求下载文件
- 更多使用示例，请看：[sample/tools/files.go](sample/tools/files.go)

```go
import(
    "context"
    "github.com/larksuite/oapi-sdk-go/core/tools"
)

// 获取文件内容
// 参数说明：
// ctx：context.Context
// url：文件的HTTP地址
// 返回值说明：
// bytes：文件内容的二进制数组
// err：错误
bytes, err := tools.DownloadFile(ctx context.Context, url string)

// 获取文件内容流，读取完文件内容后，需要关闭流
// 参数说明：
// ctx：context.Context
// url：文件的HTTP地址
// 返回值说明：
// readCloser：文件内容的二进制读取流
// err：错误
readCloser, err := tools.DownloadFileToStream(ctx context.Context, url string)

```

## 已生成SDK的业务服务

---

|业务域|版本|路径|代码示例|
|---|---|---|----|
|[用户身份验证](https://open.feishu.cn/document/ukTMukTMukTM/uETOwYjLxkDM24SM5AjN)|v1|[service/authen](service/authen)|[sample/api/authen.go](sample/api/authen.go)|
|[图片](https://open.feishu.cn/document/ukTMukTMukTM/uEDO04SM4QjLxgDN)|v4|[service/image](service/image)|[sample/api/image.go](sample/api/image.go)|
|[通讯录](https://open.feishu.cn/document/ukTMukTMukTM/uETNz4SM1MjLxUzM/v3/introduction)|v3|[service/contact](service/contact)|[sample/api/contact.go](sample/api/contact.go)|
|[日历](https://open.feishu.cn/document/ukTMukTMukTM/uETM3YjLxEzN24SMxcjN)|v4|[service/calendar](service/calendar)|[sample/api/calendar.go](sample/api/calendar.go)|
|[云空间文件](https://open.feishu.cn/document/ukTMukTMukTM/uUjM5YjL1ITO24SNykjN)|v1|[service/drive](service/drive)|[sample/api/drive.go](sample/api/drive.go)|

## License

---

- MIT

## 联系我们

---

- 飞书：[服务端SDK](https://open.feishu.cn/document/ukTMukTMukTM/uETO1YjLxkTN24SM5UjN) 页面右上角【这篇文档是否对你有帮助？】提交反馈


