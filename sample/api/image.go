package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	image "github.com/larksuite/oapi-sdk-go/service/image/v4"
	"os"
)

var imageService = image.NewService(test.GetInternalConf("online"))

func main() {
	//testUpload()
	testDownload()
}

func testUpload() {
	ctx := context.Background()
	coreCtx := core.WrapContext(ctx)
	reqCall := imageService.Images.Put(coreCtx)
	reqCall.SetImageType("message")
	f, err := os.Open("test.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	file := request.NewFile().SetContentStream(f)
	// request.NewFile().SetContent([]byte)
	reqCall.SetImage(file)
	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(result))
}

func testDownload() {
	f, err := os.Create("test_download.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	ctx := context.Background()
	coreCtx := core.WrapContext(ctx)
	reqCall := imageService.Images.Get(coreCtx)
	reqCall.SetImageKey("img_dd33673a-0f77-4bde-8ad0-xxxxxxxxx")
	reqCall.SetResponseStream(f)
	_, err = reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
}
