package main

import (
	"context"
	"fmt"
	"os"
	"time"

	client "github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/service/docx/v1"
	"github.com/larksuite/oapi-sdk-go/service/drive/v1"
)

func createDocument(client *client.Client) {
	resp, err := client.Docx.Document.Create(context.Background(), docx.NewCreateDocumentReqBuilder().
		Body(docx.NewCreateDocumentReqBodyBuilder().
			FolderToken("token").
			Title("title").
			Build()).
		Build(),
		core.WithUserAccessToken("usertoken"))

	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Code, resp.Msg, resp.RequestId(), core.Prettify(resp.Data))
}

func listBlocks(client *client.Client) {
	resp, err := client.Docx.DocumentBlock.List(context.Background(),
		docx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(100).
			Build(), core.WithUserAccessToken("u-3vEh2SpiF2WoJzYJOdiGKQ41mJrQ1hebh0G0hg.02CgW"),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.RequestId())
	fmt.Println(core.Prettify(resp))
	fmt.Println(len(resp.Data.Items))

}

func downloadFile(client *client.Client) {
	resp, err := client.Drive.File.Download(context.Background(),
		drive.NewDownloadFileReqBuilder().
			FileToken("boxcnTrRml0GB9E3NFDEyNtMeOb").
			Build(),
		core.WithUserAccessToken("u-11ETll3Kd1O8NxVwd_uVVN0hnoUAlhcbWi00kg.yyIsw"))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.RequestId())
	fmt.Println(core.Prettify(resp))

}

func listBlocksIter(client *client.Client) {
	var count = 0

	defer func() {
		fmt.Println(count)

	}()

	iter, err := client.Docx.DocumentBlock.ListDocumentBlock(context.Background(),
		docx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(1).
			Limit(3).
			Build(), core.WithUserAccessToken("u-11ETll3Kd1O8NxVwd_uVVN0hnoUAlhcbWi00kg.yyIsw"),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		hasNext, block, err := iter.Next()
		if err != nil {
			fmt.Println(err)
			return
		}

		if !hasNext {
			return
		}

		fmt.Println(core.Prettify(block))
		time.Sleep(time.Second)
		count++
	}

}

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")

	feishuClient := client.NewClient(appID, appSecret, client.WithLogLevel(core.LogLevelDebug))
	downloadFile(feishuClient)
	//listBlocks(feishuClient)
	//listBlocksIter(feishuClient)
}
