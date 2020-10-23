package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	image "github.com/larksuite/oapi-sdk-go/service/image/v4"
	"os"
)

var imageService = image.NewService(conf)

func main() {
	testUpload()
}

func testUpload() {
	ctx := context.Background()
	coreCtx := core.WarpContext(ctx)
	reqCall := imageService.Images.Put(coreCtx, request.SetTenantKey("[tenant_key]"))
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
		return
	}
	fmt.Println(tools.Prettify(result))
}
