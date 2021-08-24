# 使用应用商店应用调用服务端API

---

- [如何获取 app_access_token](https://open.feishu.cn/document/ukTMukTMukTM/uEjNz4SM2MjLxYzM) （应用商店应用）
    - 与企业自建应用相比，应用商店应用的获取 app_access_token 的流程复杂一些。
        - 需要开放平台下发的app_ticket，通过订阅事件接收。SDK已经封装了 app_ticket 事件的处理，只需要启动事件订阅服务。
        - 使用SDK调用服务端API时，如果当前还没有收到开发平台下发的 app_ticket ，会报错且向开放平台申请下发 app_ticket ，可以尽快的收到开发平台下发的 app_ticket，保证再次调用服务端API的正常。
        - 使用SDK调用服务端API时，需要使用 tenant_access_token 访问凭证时，需要 tenant_key ，来表示当前是哪个租户使用这个应用调用服务端API。
            - tenant_key，租户安装启用了这个应用，开放平台发送的服务端事件，事件内容中都含有 tenant_key。

## 使用`应用商店应用`访问 [修改用户部分信息API](https://open.feishu.cn/document/contact/v3/user/patch) 示例

- 第一步：启动事件订阅服务，用于接收`app_ticket`。

```go
package main

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	eventhttpserver "github.com/larksuite/oapi-sdk-go/event/http/native"
	"net/http"
)

var conf *config.Config

func init() {
	// 应用商店应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// EncryptKey、VerificationToken："开发者后台" -> "事件订阅" -> 事件订阅（Encrypt Key、Verification Token）
	// HelpDeskID、HelpDeskToken：https://open.feishu.cn/document/ukTMukTMukTM/ugDOyYjL4gjM24CO4IjN
	// 更多介绍请看：Github->README.zh.md->如何构建应用配置（AppSettings）
	appSettings := core.NewISVAppSettings(
		core.SetAppCredentials("AppID", "AppSecret"), // 必需
		core.SetAppEventKey("VerificationToken", "EncryptKey"), // 非必需，订阅事件、消息卡片时必需
		core.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken")) // 非必需，使用服务台API时必需

	// 当前访问的是飞书，使用默认的内存存储（app/tenant access token）、默认日志（Error级别）
	// 更多介绍请看：Github->README.zh.md->如何构建整体配置（Config）
	conf = core.NewConfig(core.DomainFeiShu, appSettings, core.SetLoggerLevel(core.LoggerLevelError))
}

func main() {
	// 启动httpServer，"开发者后台" -> "事件订阅" 请求网址 URL：https://domain/webhook/event
	eventhttpserver.Register("/webhook/event", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
```

- 第二步：在 [service](../service) 下的业务 API 或 Event，都是可以直接使用SDK。

```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
)

var conf *config.Config

func init() {
	// 应用商店应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// EncryptKey、VerificationToken："开发者后台" -> "事件订阅" -> 事件订阅（Encrypt Key、Verification Token）
	// HelpDeskID、HelpDeskToken：https://open.feishu.cn/document/ukTMukTMukTM/ugDOyYjL4gjM24CO4IjN
	// 更多介绍请看：Github->README.zh.md->如何构建应用配置（AppSettings）
	appSettings := core.NewISVAppSettings(
		core.SetAppCredentials("AppID", "AppSecret"),           // 必需
		core.SetAppEventKey("VerificationToken", "EncryptKey"), // 非必需，订阅事件、消息卡片时必需
		core.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken")) // 非必需，使用服务台API时必需

	// 当前访问的是飞书，使用默认的内存存储（app/tenant access token）、默认日志（Error级别）
	// 更多介绍请看：Github->README.zh.md->如何构建整体配置（Config）
	conf = core.NewConfig(core.DomainFeiShu, appSettings, core.SetLoggerLevel(core.LoggerLevelError))
}

func main() {
	service := contact.NewService(conf)
	coreCtx := core.WrapContext(context.Background())
	body := &contact.User{}
	body.Name = "rename"
	// 由于这是一个PATCH请求，需要告之更新哪些字段
	body.ForceSendFields = append(body.ForceSendFields, "Name")
	// 构建请求 && 设置租户标识（tenant_key）
	reqCall := service.Users.Patch(coreCtx, body, request.SetTenantKey("tenant_key"))
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
## 使用`应用商店应用`访问 [发送文本消息API](https://open.feishu.cn/document/ukTMukTMukTM/uUjNz4SN2MjL1YzM) 示例
  
- 第一步：启动事件订阅服务，用于接收`app_ticket`。
  - 同上
    
- 第二步：调用服务端接口，有些老版接口，没有直接可以使用的SDK，可以使用`原生`模式。

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
	"github.com/larksuite/oapi-sdk-go/core/tools"
)

var conf *config.Config

func init() {
	// 应用商店应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// EncryptKey、VerificationToken："开发者后台" -> "事件订阅" -> 事件订阅（Encrypt Key、Verification Token）
	// HelpDeskID、HelpDeskToken：https://open.feishu.cn/document/ukTMukTMukTM/ugDOyYjL4gjM24CO4IjN
	// 更多介绍请看：Github->README.zh.md->如何构建应用配置（AppSettings）
	appSettings := core.NewISVAppSettings(
		core.SetAppCredentials("AppID", "AppSecret"),           // 必需
		core.SetAppEventKey("VerificationToken", "EncryptKey"), // 非必需，订阅事件、消息卡片时必需
		core.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken")) // 非必需，使用服务台API时必需

	// 当前访问的是飞书，使用默认的内存存储（app/tenant access token）、默认日志（Error级别）
	// 更多介绍请看：Github->README.zh.md->如何构建整体配置（Config）
	conf = core.NewConfig(core.DomainFeiShu, appSettings, core.SetLoggerLevel(core.LoggerLevelError))
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
	// 构建请求&&设置企业标识（tenant_key）
	req := request.NewRequestWithNative("/open-apis/message/v4/send", "POST", request.AccessTokenTypeTenant,
		body, &ret, request.SetTenantKey("Tenant key"))
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


