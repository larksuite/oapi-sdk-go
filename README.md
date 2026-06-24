# Feishu OpenPlatform Server SDK

[English](./README.md) | [Simplified Chinese](./README.zh.md)

Feishu Open Platform offers a series of server-side atomic APIs to achieve diverse functionalities. However, actual coding requires additional work, such as obtaining and maintaining access tokens, encrypting and decrypting data, and verifying request signatures. Furthermore, the lack of semantic descriptions for function calls and type system support can increase coding burdens.

To address these issues, Feishu Open Platform has developed the Open Interface SDK, which incorporates these lengthy logic processes, provides a comprehensive type system, and offers a semantic programming interface to improve the coding experience.

## Introduction Documents

- [Preparations before development (Install SDK)](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/preparations)
- [Calling Server-side APIs](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/calling-server-side-apis)
- [Handle Events](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/handle-events)
- [Handle Card Callbacks](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/golang-sdk-guide/handle-callback)
- [SDK FAQs](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/server-side-sdk/faq)

## Channel Module

The SDK provides a `Channel` module built on top of WebSocket and the API Client. It encapsulates event listening, message normalization, streaming replies, and media uploads, allowing developers to focus purely on business logic.

- [Channel Module Documentation (English)](./doc/channel.md)
- [Channel Module Documentation (Chinese)](./doc/channel.zh.md)

## One-Click App Registration

The SDK provides `registration.RegisterApp` for one-click app creation based on OAuth 2.0 Device Authorization Grant (RFC 8628). It returns a verification URL that users can open in Feishu/Lark or render as a QR code. After authorization, the app is created automatically and the SDK returns the `App ID` and `App Secret`.

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

### Custom Scopes, Events, Callbacks, And Updating An Existing App

When creating an app, use `Options.Addons` to incrementally request scopes, event subscriptions, and callbacks on top of the platform base template. The values are pre-filled into the confirmation page after the user opens the QR code URL and take effect after confirmation. `Options.CreateOnly=true` only allows creating a new app. `Options.AppID` starts the update flow for an existing app.

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

Notes: `Addons` is additive only and cannot remove config from the base template. The SDK validates shape and non-empty strings, but does not validate whether scope, event, or callback names exist.

### `registration.RegisterApp` Parameters

| Parameter | Description | Type | Required | Default |
| ---- | ---- | ---- | ---- | ---- |
| `ctx` | Controls timeout and cancellation for the registration flow; canceling the `context` stops polling. | `context.Context` | Yes | - |
| `Options.Source` | Source identifier appended to the QR URL as `go-sdk/{source}`. | `string` | No | `go-sdk` |
| `Options.Domain` | Custom Feishu accounts domain. A full base URL such as `https://accounts.feishu.cn` is supported. | `string` | No | `https://accounts.feishu.cn` |
| `Options.LarkDomain` | Custom Lark accounts domain used when `tenant_brand=lark` is detected. | `string` | No | `https://accounts.larksuite.com` |
| `Options.AppPreset` | Pre-filled app creation values; users can still edit them on the page. | `*registration.AppPreset` | No | - |
| `Options.AppPreset.Avatar` | App avatar URLs, 1-6 entries; first entry is selected by default. Pass raw URLs and the SDK encodes them. Page-side display rules are handled by the app creation page. | `[]string` | No | - |
| `Options.AppPreset.Name` | App name with `{user}` placeholder support; pass raw value and the SDK encodes it. | `string` | No | - |
| `Options.AppPreset.Desc` | App description with `{user}` placeholder support; pass raw value and the SDK encodes it. | `string` | No | - |
| `Options.Addons` | Incremental scopes, events, and callbacks pre-filled into the confirmation page. | `*registration.AppAddons` | No | - |
| `Options.Addons.Scopes.Tenant` | App-identity scopes, for example `im:message:send_as_bot`. | `[]string` | No | - |
| `Options.Addons.Scopes.User` | User-identity scopes, for example `calendar:calendar:read`. | `[]string` | No | - |
| `Options.Addons.Events.Items.Tenant` | App-identity events, for example `im.message.receive_v1`. | `[]string` | No | - |
| `Options.Addons.Events.Items.User` | User-identity events, for example `calendar.calendar.event.changed_v4`. | `[]string` | No | - |
| `Options.Addons.Callbacks.Items` | Callback names, for example `card.action.trigger`. | `[]string` | No | - |
| `Options.CreateOnly` | When `true`, the landing page only allows creating a new app. When used together with `Options.AppID`, the page gives create-new-app flow precedence. | `bool` | No | `false` |
| `Options.AppID` | Existing app ID carried as `clientID` in the QR URL for the update flow. | `string` | No | - |
| `Options.OnQRCode` | Callback invoked when the verification URL is ready. The callback receives `{ URL, ExpireIn }`. | `func(info *registration.QRCodeInfo)` | Yes | - |
| `Options.OnStatusChange` | Callback for polling status changes. `Status` can be `polling`, `slow_down`, or `domain_switched`. | `func(info *registration.StatusChangeInfo)` | No | - |

### Return Value

| Field | Type | Description |
| ---- | ---- | ---- |
| `ClientID` | `string` | App ID |
| `ClientSecret` | `string` | App Secret |
| `UserInfo` | `*registration.UserInfo` | Scanning user info |
| `UserInfo.OpenID` | `string` | User `open_id` |
| `UserInfo.TenantBrand` | `string` | `"feishu"` or `"lark"` |

### Error Handling

Returned errors usually expose `Code` and `Description` through `registration.RegisterAppError`. More specific types include `registration.AccessDeniedError` and `registration.ExpiredError`.

| `Code` | Description |
| ---- | ---- |
| `access_denied` | User denied authorization |
| `expired_token` | QR code expired or polling timed out |
| `invalid_response` | Response is empty or missing required fields |

## Extended Examples

We also provide common API composition examples and business scenario examples based on the SDK, such as:

* Messages
    * [Send file messages](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/im/send_file.go)
    * [Send image messages](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/im/send_image.go)
* Contacts
    * [Get all users under a department](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/contact/list_user_by_department.go)
* Base
    * [Create a Base app and add tables](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/base/create_app_and_tables.go)
* Sheets
    * [Copy and paste a cell range](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/sheets/copy_and_paste_by_range.go)
    * [Download all media in a cell range](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/composite_api/sheets/download_media_by_range.go)
* Tutorials
    * [Robot alert group automation](https://github.com/larksuite/oapi-sdk-go-demo/blob/main/quick_start/robot) ([development tutorial](https://open.feishu.cn/document/home/message-development-tutorial/introduction))

For more examples, see https://github.com/larksuite/oapi-sdk-go-demo

## Community

[Join the support group](https://applink.feishu.cn/client/chat/chatter/add_by_link?link_token=985nb30c-787a-4fbb-904d-2cf945534078)

## License

MIT
