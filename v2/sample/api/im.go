package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/service/im/v1"
	"io/ioutil"
	"os"
)

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	ctx := context.Background()
	messageCreate(ctx, larkApp)
	fileCreate(ctx, larkApp)
	fileDownload(ctx, larkApp)
}

func messageCreate(ctx context.Context, larkApp *lark.App) {
	messageText := &lark.MessageText{Text: "Tom test content"}
	content, err := messageText.JSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	messageCreateResp, err := im.New(larkApp).Messages.Create(ctx, &im.MessageCreateReq{
		ReceiveIdType: lark.StringPtr("user_id"),
		Body: &im.MessageCreateReqBody{
			ReceiveId: lark.StringPtr("77bbc392"),
			MsgType:   lark.StringPtr("text"),
			Content:   lark.StringPtr(content),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", messageCreateResp.RequestId())
	if messageCreateResp.Code != 0 {
		fmt.Println(messageCreateResp.CodeError)
		return
	}
	fmt.Println(lark.Prettify(messageCreateResp.Data))
	fmt.Println()
	fmt.Println()
}

func fileCreate(ctx context.Context, larkApp *lark.App) {
	pdf, err := os.Open("test.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pdf.Close()

	fileCreateResp, err := im.New(larkApp).Files.Create(ctx, &im.FileCreateReq{
		Body: &im.FileCreateReqBody{
			FileType: lark.StringPtr("pdf"),
			FileName: lark.StringPtr("测试.pdf"),
			File:     pdf,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", fileCreateResp.RequestId())
	if fileCreateResp.Code != 0 {
		fmt.Println(fileCreateResp.CodeError)
		return
	}
	fmt.Println(lark.Prettify(fileCreateResp))
	fmt.Println()
	fmt.Println()
}

func fileDownload(ctx context.Context, larkApp *lark.App) {
	fileGetResp, err := im.New(larkApp).Files.Get(ctx, &im.FileGetReq{
		FileKey: "file_v2_62ac7c6e-de7e-464f-ac33-f1c28f94169g",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", fileGetResp.RequestId())
	if fileGetResp.Code != 0 {
		fmt.Println(fileGetResp.CodeError)
		return
	}
	fmt.Printf("file name:%s \n", fileGetResp.FileName)

	bs, err := ioutil.ReadAll(fileGetResp.File)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("test_download_v2.pdf", bs, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}
