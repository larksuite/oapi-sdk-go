package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	im "github.com/larksuite/oapi-sdk-go/v2/service/im/v1"
	"io/ioutil"
	"os"
)

func messageCreate(ctx context.Context, larkApp *lark.App) {
	messageCreateResp, err := im.New(larkApp).Messages.Create(ctx, &im.MessageCreateReq{
		ReceiveIdType: im.ReceiveIdTypeUserId.Ptr(),
		Body: im.MessageCreateReqBody{
			ReceiveId: lark.StringPtr("77bbc392"),
			MsgType:   lark.StringPtr("text"),
			Content:   lark.StringPtr(`{"text":"Tom test content"}`),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("request id: %s \n", messageCreateResp.RequestId())
	if messageCreateResp.Code != 0 {
		panic(messageCreateResp.CodeError.Error())
	}
	fmt.Println(lark.Prettify(messageCreateResp.Data))
}

func fileCreate(ctx context.Context, larkApp *lark.App) {
	pdf, err := os.Open("test.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pdf.Close()

	fileCreateResp, err := im.New(larkApp).Files.Create(ctx, &im.FileCreateReq{
		FileType: lark.StringPtr("pdf"),
		FileName: lark.StringPtr("测试.pdf"),
		File:     pdf,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("request id: %s \n", fileCreateResp.RequestId())
	if fileCreateResp.Code != 0 {
		panic(fileCreateResp.CodeError.Error())
	}
	fmt.Println(lark.Prettify(fileCreateResp))
}

func fileDownload(ctx context.Context, larkApp *lark.App) {
	fileGetResp, err := im.New(larkApp).Files.Get(ctx, &im.FileGetReq{
		FileKey: "file_v2_0c24c84e-819d-417f-b42f-978cf1b50aag",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("request id: %s \n", fileGetResp.RequestId())
	if fileGetResp.Code != 0 {
		panic(fileGetResp.CodeError.Error())
	}
	fmt.Printf("file name:%s \n", fileGetResp.FileName)

	bs, err := ioutil.ReadAll(fileGetResp.File)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("test_download_v2.pdf", bs, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	larkApp := lark.NewApp(lark.DomainFeiShu, lark.WithAppCredential(appID, appSecret),
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	ctx := context.Background()
	messageCreate(ctx, larkApp)
	fileCreate(ctx, larkApp)
	fileDownload(ctx, larkApp)
}
