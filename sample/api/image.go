package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	image "github.com/larksuite/oapi-sdk-go/service/image/v4"
	"os"
)

var appConf *lark.AppConfig

func main() {
	appConf = lark.NewInternalAppConfig(lark.DomainFeiShu, lark.SetAppCredentials("AppID", "AppSecret"), // 必需
		lark.SetAppEventKey("VerificationToken", "EncryptKey"),     // 非必需，订阅事件、消息卡片时必需
		lark.SetHelpDeskCredentials("HelpDeskID", "HelpDeskToken")) // 非必需，使用服务台API时必需)
	appConf.SetLogLevel(lark.LogLevelDebug)
	testUpload()
	testDownload()
}

func testUpload() {
	ctx := context.Background()
	coreCtx := lark.WrapContext(ctx)
	reqCall := image.NewService(appConf).Images.Put(coreCtx)
	reqCall.SetImageType("message")
	f, err := os.Open("test.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	file := lark.NewFormDataFile().SetContentStream(f)
	// lark.NewFormDataFile().SetContent([]byte)
	reqCall.SetImage(file)
	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(result))
}

func testDownload() {
	f, err := os.Create("test_download.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	ctx := context.Background()
	coreCtx := lark.WrapContext(ctx)
	reqCall := image.NewService(appConf).Images.Get(coreCtx)
	reqCall.SetImageKey("img_800c6035-7db8-4844-bc85-01a74d6e5cag")
	reqCall.SetResponseStream(f)
	_, err = reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
}
