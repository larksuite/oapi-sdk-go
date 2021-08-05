package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
)

func main() {
	// 应用商店应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// EncryptKey、VerificationToken："开发者后台" -> "事件订阅" -> 事件订阅（Encrypt Key、Verification Token）
	// HelpDeskID、HelpDeskToken：https://open.feishu.cn/document/ukTMukTMukTM/ugDOyYjL4gjM24CO4IjN
	// 当前访问的是飞书，使用默认的内存存储（app/tenant access token）、默认日志（Error级别）
	// 更多介绍请看：Github->README.zh.md->如何构建整体配置（Config）
	conf := lark.NewInternalAppConfig(lark.DomainFeiShu, lark.SetAppCredentials("AppID", "AppSecret"), // 必需
		lark.SetAppEventKey("VerificationToken", "EncryptKey"),     // 非必需，订阅事件、消息卡片时必需
		lark.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken")) // 非必需，使用服务台API时必需)
	conf.SetLogLevel(lark.LogLevelDebug)
	// 请求发送消息的结果
	ret := make(map[string]interface{})
	// 构建请求
	req := lark.NewAPIRequest("/open-apis/helpdesk/v1/tickets/6971250929135779860", "GET",
		lark.AccessTokenTypeTenant, nil, &ret,
		lark.NeedHelpDeskAuth(), // 服务台 API，需要 HelpDeskAuth
	)
	// 请求的上下文
	coreCtx := lark.WrapContext(context.Background())
	// 发送请求
	err := lark.SendAPIRequest(coreCtx, conf, req)
	// 打印请求的RequestID
	fmt.Println(coreCtx.GetRequestID())
	// 打印请求的响应状态吗
	fmt.Println(coreCtx.GetHTTPStatusCode())
	// 请求的error处理
	if err != nil {
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		fmt.Println(lark.Prettify(err))
		return
	}
	// 打印请求的结果
	fmt.Println(lark.Prettify(ret))
}
