# 飞书开放接口 SDK

[English](./README.md) | [简体中文](./README.zh.md)

飞书开放平台提供了一系列服务端原子 API，用于实现多种业务能力。但在实际编码中，开发者通常还需要处理 access token 获取与维护、数据加解密、请求签名校验等通用逻辑；同时，缺少语义化调用方式和类型系统支持也会增加编码负担。

为了解决这些问题，飞书开放平台提供了开放接口 SDK。SDK 封装了冗长的通用逻辑，提供完善的类型系统和语义化编程接口，帮助开发者提升编码体验。

## 介绍文档

- [开发前准备（安装 SDK）](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/preparations)
- [调用服务端 API](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/calling-server-side-apis)
- [处理事件订阅](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/handle-events)
- [处理卡片回调](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/handle-callback)
- [常见问题](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/faq)

## 高级封装 Channel 模块

SDK 提供了一个基于 WebSocket 和 API Client 封装的 Channel 模块。它将飞书机器人接入过程中的事件监听、消息归一化、发送流式回复、上传媒体等操作进行了高层封装，让开发者能更专注业务逻辑。

- [Channel 模块文档（中文）](./doc/channel.zh.md)
- [Channel Module Documentation (English)](./doc/channel.md)

## 一键创建应用

SDK 提供了 `registration.RegisterApp` 方法，基于 OAuth 2.0 Device Authorization Grant（RFC 8628）协议实现一键创建应用。调用后会返回一个验证链接，用户在飞书/Lark 中打开链接或扫码完成授权后，即可自动注册应用并获取 `App ID` 与 `App Secret`，无需手动前往开发者后台创建。

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/scene/registration"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	result, err := registration.RegisterApp(ctx, &registration.Options{
		OnQRCode: func(info *registration.QRCodeInfo) {
			fmt.Printf("open or scan this url: %s\n", info.URL)
			fmt.Printf("the link expires in %d seconds\n", info.ExpireIn)
		},
		OnStatusChange: func(info *registration.StatusChangeInfo) {
			// status: polling | slow_down | domain_switched
			fmt.Printf("registration status: %s", info.Status)
			if info.Interval > 0 {
				fmt.Printf(", next poll after %d seconds", info.Interval)
			}
			fmt.Println()
		},
	})
	if err != nil {
		var regErr *registration.RegisterAppError
		if errors.As(err, &regErr) {
			fmt.Printf("register app failed: code=%s, description=%s\n", regErr.Code, regErr.Description)
			return
		}
		panic(err)
	}

	fmt.Println("App ID:", result.ClientID)
	fmt.Println("App Secret:", result.ClientSecret)

	client := lark.NewClient(result.ClientID, result.ClientSecret)
	_ = client
}
```

### 自定义权限、事件、回调与更新已有应用

创建应用时，可通过 `Options.Addons` 在平台基础模板上增量申请权限、事件订阅和回调。这些配置会预填到用户扫码后的确认页中，用户确认后生效。`Options.CreateOnly=true` 时，页面只允许创建新应用；传入 `Options.AppID` 时，页面会进入更新已有应用配置的确认流程。

```go
_, err := registration.RegisterApp(ctx, &registration.Options{
	Addons: &registration.AppAddons{
		Scopes: registration.AppAddonsScopes{
			Tenant: []string{"im:message:send_as_bot"},
			User:   []string{"calendar:calendar:read"},
		},
		Events: registration.AppAddonsEvents{
			Items: registration.AppAddonsEventItems{
				Tenant: []string{"im.message.receive_v1"},
			},
		},
		Callbacks: registration.AppAddonsCallbacks{
			Items: []string{"card.action.trigger"},
		},
	},
	CreateOnly: true,
	OnQRCode: func(info *registration.QRCodeInfo) {
		fmt.Println(info.URL)
	},
})
if err != nil {
	panic(err)
}

_, err = registration.RegisterApp(ctx, &registration.Options{
	AppID: "cli_xxx",
	Addons: &registration.AppAddons{
		Scopes: registration.AppAddonsScopes{
			Tenant: []string{"drive:drive.metadata:readonly"},
		},
	},
	OnQRCode: func(info *registration.QRCodeInfo) {
		fmt.Println(info.URL)
	},
})
if err != nil {
	panic(err)
}
```

注意：`Addons` 仅支持增量叠加，不支持删除基础模板中的配置；SDK 只校验数据形状和非空字符串，不校验权限点、事件或回调名称是否存在。

### `registration.RegisterApp` 参数

| 参数 | 描述 | 类型 | 必填 | 默认值 |
| ---- | ---- | ---- | ---- | ---- |
| `ctx` | 控制注册流程的超时与取消；取消 `context` 会终止轮询。 | `context.Context` | 是 | - |
| `Options.Source` | 来源标识，会拼入二维码 URL 的 `source` 参数，格式为 `go-sdk/{source}`。 | `string` | 否 | `go-sdk` |
| `Options.Domain` | 自定义飞书认证域名，支持传完整前缀，如 `https://accounts.feishu.cn`。 | `string` | 否 | `https://accounts.feishu.cn` |
| `Options.LarkDomain` | 自定义 Lark 认证域名；检测到 `tenant_brand=lark` 时自动切换。 | `string` | 否 | `https://accounts.larksuite.com` |
| `Options.AppPreset` | 预设应用信息，仅用于初始化创建页；用户仍可在页面修改，最终以页面提交为准。 | `*registration.AppPreset` | 否 | - |
| `Options.AppPreset.Avatar` | 应用头像 URL，支持 1-6 个；第一个默认选中。传原始 URL，SDK 会编码。头像展示、图片可访问性、GIF 取帧等由创建页处理。 | `[]string` | 否 | - |
| `Options.AppPreset.Name` | 应用名称，支持 `{user}` 占位符；传原始值，SDK 会编码。 | `string` | 否 | - |
| `Options.AppPreset.Desc` | 应用描述，支持 `{user}` 占位符；传原始值，SDK 会编码。 | `string` | 否 | - |
| `Options.Addons` | 增量权限、事件、回调配置，预填到扫码后的确认页，用户确认后生效。 | `*registration.AppAddons` | 否 | - |
| `Options.Addons.Scopes.Tenant` | 应用身份权限列表，如 `im:message:send_as_bot`。 | `[]string` | 否 | - |
| `Options.Addons.Scopes.User` | 用户身份权限列表，如 `calendar:calendar:read`。 | `[]string` | 否 | - |
| `Options.Addons.Events.Items.Tenant` | 应用身份事件列表，如 `im.message.receive_v1`。 | `[]string` | 否 | - |
| `Options.Addons.Events.Items.User` | 用户身份事件列表，如 `calendar.calendar.event.changed_v4`。 | `[]string` | 否 | - |
| `Options.Addons.Callbacks.Items` | 回调列表，如 `card.action.trigger`。 | `[]string` | 否 | - |
| `Options.CreateOnly` | 为 `true` 时落地页仅允许创建新应用；与 `Options.AppID` 同时传入时，页面侧优先走新建流程。 | `bool` | 否 | `false` |
| `Options.AppID` | 已有应用的 App ID；传入后会作为二维码 URL 的 `clientID` 参数，用于更新已有应用配置。 | `string` | 否 | - |
| `Options.OnQRCode` | 验证链接就绪时的回调，参数为 `{ URL, ExpireIn }`。可直接展示链接，或将其渲染为二维码供用户扫码。 | `func(info *registration.QRCodeInfo)` | 是 | - |
| `Options.OnStatusChange` | 轮询状态变化回调，参数为 `{ Status, Interval }`。`Status` 可能为 `polling`、`slow_down`、`domain_switched`。 | `func(info *registration.StatusChangeInfo)` | 否 | - |

### 返回值

| 字段 | 类型 | 描述 |
| ---- | ---- | ---- |
| `ClientID` | `string` | 应用的 `App ID` |
| `ClientSecret` | `string` | 应用的 `App Secret` |
| `UserInfo` | `*registration.UserInfo` | 扫码用户信息 |
| `UserInfo.OpenID` | `string` | 扫码用户的 `open_id` |
| `UserInfo.TenantBrand` | `string` | `"feishu"` 或 `"lark"` |

### 错误处理

返回的错误通常可以通过 `errors.As(err, &registration.RegisterAppError)` 获取 `Code` 与 `Description` 字段。更具体的错误类型还包括 `registration.AccessDeniedError` 和 `registration.ExpiredError`。

| `Code` | 描述 |
| ---- | ---- |
| `access_denied` | 用户拒绝授权 |
| `expired_token` | 二维码过期或轮询超时 |
| `invalid_response` | 接口返回缺少必要字段或响应为空 |

## 扩展示例

我们还基于 SDK 封装了常用的 API 组合调用及业务场景示例，如：

* 消息
    * [发送文件消息](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/im/send_file.go)
    * [发送图片消息](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/im/send_image.go)
* 通讯录
    * [获取部门下所有用户列表](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/contact/list_user_by_department.go)
* 多维表格
    * [创建多维表格同时添加数据表](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/base/create_app_and_tables.go)
* 电子表格
    * [复制粘贴某个范围的单元格数据](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/sheets/copy_and_paste_by_range.go)
    * [下载指定范围单元格的所有素材列表](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/sheets/download_media_by_range.go)
* 教程
    * [机器人自动拉群报警](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/quick_start/robot)（[开发教程](https://open.feishu.cn/document/home/message-development-tutorial/introduction)）

更多示例可参考：https://github.com/larksuite/oapi-sdk-go-demo

## 加入交流互助群

[单击加入交流互助](https://applink.feishu.cn/client/chat/chatter/add_by_link?link_token=985nb30c-787a-4fbb-904d-2cf945534078)

## License

MIT
