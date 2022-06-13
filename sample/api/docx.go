package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/feishu/oapi-sdk-go"
	"github.com/feishu/oapi-sdk-go/core"
	"github.com/feishu/oapi-sdk-go/service/docx/v1"
)

func createDocument(client *client.Client) {
	resp, err := client.Docx.Documents.Create(context.Background(), docx.NewCreateDocumentReqBuilder().
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
	resp, err := client.Docx.Blocks.List(context.Background(),
		docx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(100).
			Build(), core.WithUserAccessToken("u-JDboiwm9RnJbtNc1gdJ0Qd"),
	)

	if err != nil {
		fmt.Println(core.Prettify(err))
		return
	}

	fmt.Println(resp.RequestId())
	fmt.Println(core.Prettify(resp))
	fmt.Println(len(resp.Data.Items))

}

func listBlocksIter(client *client.Client) {
	var count = 0

	defer func() {
		fmt.Println(count)

	}()

	iter, err := client.Docx.Blocks.ListDocumentBlock(context.Background(),
		docx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(2).
			Limit(100).
			Build(), core.WithUserAccessToken("u-JDboiwm9RnJbtNc1gdJ0Qd"),
	)

	if err != nil {
		fmt.Println(core.Prettify(err))
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

	feishuClient := client.NewClient(appID, appSecret)

	listBlocks(feishuClient)
	listBlocksIter(feishuClient)
}
