package main

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/service/docx/v1"
)

var userToken = "u-WsRMhBN7CmTGoBL2mqZ7Ha"

func createDocument(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)

	resp, err := docxService.Documents.Create(ctx, &docx.DocumentCreateReq{Body: &docx.DocumentCreateReqBody{
		FolderToken: lark.StringPtr("fldcniHf40Vcv1DoEc8SXeuA0Zd"),
		Title:       lark.StringPtr("document1"),
	}}, lark.WithUserAccessToken(userToken))

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func getDocument(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)

	resp, err := docxService.Documents.Get(ctx, &docx.DocumentGetReq{DocumentId: "doxcneIVZ4VJcVDNWqTlSQ4dzrd"})

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func getRawDocument(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.Documents.RawContent(ctx, &docx.DocumentRawContentReq{
		DocumentId: "doxcneIVZ4VJcVDNWqTlSQ4dzrd",
		Lang:       lark.IntPtr(0)})

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

var DocumentId = "doxcn0stWG7Zb9ItETHUVHv6fsg"

func listBlocks(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.Blocks.List(ctx, &docx.DocumentBlockListReq{
		DocumentId: DocumentId,
		PageSize:   lark.IntPtr(10),
	}, lark.WithUserAccessToken(userToken))

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func createBlocks(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.DocumentBlockChildren.Create(ctx, &docx.DocumentBlockChildrenCreateReq{
		DocumentId: DocumentId,
		BlockId:    "doxcnvDRZwWqR4H2tLvDWi3JHxf",
		//ClientToken: lark.StringPtr("ssss"),
		Body: &docx.DocumentBlockChildrenCreateReqBody{
			Children: []*docx.Block{&docx.Block{
				BlockType: lark.IntPtr(2),
				Text: &docx.Text{
					Elements: []*docx.TextElement{&docx.TextElement{
						TextRun: &docx.TextRun{
							Content: lark.StringPtr("插入字块1"),
							TextElementStyle: &docx.TextElementStyle{
								BackgroundColor: lark.IntPtr(14),
								TextColor:       lark.IntPtr(5),
							},
						},
					}, &docx.TextElement{
						TextRun: &docx.TextRun{
							Content: lark.StringPtr("插入字块2"),
							TextElementStyle: &docx.TextElementStyle{
								BackgroundColor: lark.IntPtr(14),
								TextColor:       lark.IntPtr(5),
								Bold:            lark.BoolPtr(true),
							},
						},
					}},
					Style: &docx.TextStyle{},
				},
			}},
			Index: lark.IntPtr(0),
		},
	}, lark.WithUserAccessToken(userToken))

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func updateBlock(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.Blocks.Patch(ctx, &docx.DocumentBlockPatchReq{
		DocumentId: DocumentId,
		BlockId:    "doxcnIWUK8qWm2cMWyWwnHNd9Gc",
		UpdateBlockRequest: &docx.UpdateBlockRequest{
			UpdateTextElements: &docx.UpdateTextElementsRequest{
				Elements: []*docx.TextElement{&docx.TextElement{
					TextRun: &docx.TextRun{
						Content: lark.StringPtr("ssssssss"),
					},
				}},
			},
		},
	}, lark.WithUserAccessToken(userToken))

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func batchUpdateBlock(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.Blocks.BatchUpdate(ctx, &docx.DocumentBlockBatchUpdateReq{
		DocumentId: DocumentId,
		Body:       &docx.DocumentBlockBatchUpdateReqBody{Requests: []*docx.UpdateBlockRequest{}},
	})

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func getAllSubBlock(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.DocumentBlockChildren.Get(ctx, &docx.DocumentBlockChildrenGetReq{
		DocumentId:         DocumentId,
		BlockId:            "doxcnvDRZwWqR4H2tLvDWi3JHxf",
		DocumentRevisionId: nil,
		PageToken:          nil,
		PageSize:           nil,
		UserIdType:         nil,
	}, lark.WithUserAccessToken(userToken))

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func getBlock(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.Blocks.Get(ctx, &docx.DocumentBlockGetReq{
		DocumentId:         DocumentId,
		BlockId:            "doxcnkMWa6MGsEY02wfp22Us7ef",
		DocumentRevisionId: nil,
		UserIdType:         nil,
	}, lark.WithUserAccessToken(userToken))

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func batchDelBlock(larkApp *lark.App) {
	ctx := context.Background()
	docxService := docx.New(larkApp)
	resp, err := docxService.DocumentBlockChildren.BatchDelete(ctx, &docx.DocumentBlockChildrenBatchDeleteReq{
		DocumentId: DocumentId,
		BlockId:    "doxcnvDRZwWqR4H2tLvDWi3JHxf",
		Body: &docx.DocumentBlockChildrenBatchDeleteReqBody{
			StartIndex: lark.IntPtr(0),
			EndIndex:   lark.IntPtr(1),
		},
	}, lark.WithUserAccessToken(userToken))

	if err != nil {
		panic(err)
	}

	if resp.Code != 0 {
		fmt.Println(resp.Code, resp.Msg)
		return
	}

	fmt.Println(lark.Prettify(resp.Data))
}

func main() {

	var appID, appSecret = "xxx", "xxx"
	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	// 创建文档
	//createDocument(larkApp)

	// 获取文档基本信息
	//getDocument(larkApp)

	//// 获取纯文本信息
	//getRawDocument(larkApp)

	//创建块
	//createBlocks(larkApp)

	// 获取文档block列表
	listBlocks(larkApp)

	// 更新块
	//updateBlock(larkApp)

	// 批量更新block
	//batchUpdateBlock(larkApp)

	// 获取所有子块
	//getAllSubBlock(larkApp)

	// 获取块
	//getBlock(larkApp)

	//批量删除块
	//batchDelBlock(larkApp)

}
