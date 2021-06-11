package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	drive_explorer "github.com/larksuite/oapi-sdk-go/service/drive_explorer/v2"
)

func main() {
	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）。
	// 更多介绍请看：Github->README.zh.md->高级使用->如何构建应用配置（AppSettings）
	appSetting := config.NewInternalAppSettings("AppID", "AppSecret", "VerificationToken", "EncryptKey")

	// 当前访问的是飞书，使用默认的内存存储（app/tenant access token）、默认日志（Debug级别）
	// 更多介绍请看：Github->README.zh.md->高级使用->如何构建整体配置（Config）
	conf = config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelDebug)

	service := drive_explorer.NewService(conf)
	coreCtx := core.WrapContext(context.Background())
	var optFns []request.OptFn
	optFns = append(optFns, request.SetUserAccessToken("u-xxxx"))
	// body params
	body := &drive_explorer.FileCreateReqBody{}
	body.Title = "33"
	body.Type = "444"
	reqCall := service.Files.Create(coreCtx, body, optFns...)
	// path params
	reqCall.SetFolderToken("ddddddd")
	result, err := reqCall.Do()
	fmt.Printf("request id:%s\n", coreCtx.GetRequestID())
	fmt.Printf("HTTP status code:%d\n", coreCtx.GetHTTPStatusCode())
	fmt.Printf("HTTP response header:%s\n", coreCtx.GetHeader())
	if err != nil {
		realErr := err.(*response.Error)
		fmt.Printf("err code:%d\n", realErr.Code)
		fmt.Printf("err msg:%s\n", realErr.Msg)
		fmt.Printf("err:%v\n", realErr)
		return
	}
	fmt.Println(result)
	//fmt.Printf("result:%v\n", tools.Prettify(result))

}
