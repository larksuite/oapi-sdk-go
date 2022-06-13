package main

import (
	"context"
	"fmt"

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
	fmt.Println(resp.Code, resp.Msg, core.Prettify(resp.Data))
}

func listBlocks(client *client.Client) {
	resp, err := client.Docx.Blocks.List(context.Background(),
		docx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(1).
			Build(), core.WithUserAccessToken("u-kFK7mQdQasTbiosC18boUc"),
	)

	if err != nil {
		fmt.Println(core.Prettify(err))
		return
	}

	fmt.Println(core.Prettify(resp))
}

//
//func listBlocksIter() {
//	iter, err := client.Docx.Blocks.ListDocumentBlock(context.Background(),
//		docx.NewListDocumentBlockReqBuilder().
//			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
//			PageSize(2).
//			Build(), core.WithUserAccessToken("u-zwbYaTxHGGHxQ9BAVIAO5g"),
//	)
//
//	if err != nil {
//		fmt.Println(core.Prettify(err))
//		return
//	}
//
//	for {
//		if iter.HasNext() {
//			block, err := iter.Next()
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//			fmt.Println(core.Prettify(block))
//
//		} else {
//			break
//		}
//	}
//
//}

func main() {
	feishuClient := client.NewClient("cli_a1eccc36c278900d", "0PhrmTxRd7q6cqzVKx25tgvlObXNmbqD")

	listBlocks(feishuClient)
	//listBlocksIter()
}
