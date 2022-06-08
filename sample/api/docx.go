package main

import (
	"context"
	"fmt"

	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	docx "github.com/larksuite/oapi-sdk-go/service/docx/v1"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var docxService = docx.NewService(configs.TestConfig(core.DomainFeiShu))

func main() {
	//createDocument()
	//getDocument()
	//getRawDocument()
	//createBlock()
	listBlocks()
	//listSubBlocks()

	//getBlock()
	//batchDelBlock()

}
func createDocument() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.Documents.Create(coreCtx, &docx.DocumentCreateReqBody{
		FolderToken: "fldcniHf40Vcv1DoEc8SXeuA0Zd",
		Title:       "documentv1",
	}, request.SetUserAccessToken("u-GQze1ue1QXC650Z3NYBsga"))

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func getDocument() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.Documents.Get(coreCtx, request.SetUserAccessToken("u-GQze1ue1QXC650Z3NYBsga"))
	reqCall.SetDocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func getRawDocument() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.Documents.RawContent(coreCtx, request.SetUserAccessToken("u-G4p3fYOXuJqkwyNOwSDG5g"))
	reqCall.SetDocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}
func createBlock() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.DocumentBlockChildrens.Create(coreCtx, &docx.DocumentBlockChildrenCreateReqBody{
		Children: []*docx.Block{&docx.Block{
			BlockType: 2,
			Text: &docx.Text{Elements: []*docx.TextElement{
				&docx.TextElement{
					TextRun: &docx.TextRun{
						Content:          "插入v1-1块",
						TextElementStyle: &docx.TextElementStyle{BackgroundColor: 14, TextColor: 5}}},

				&docx.TextElement{
					TextRun: &docx.TextRun{
						Content:          "插入v1-2块",
						TextElementStyle: &docx.TextElementStyle{BackgroundColor: 14, TextColor: 5}}}}},
		}},
		Index:           0,
		ForceSendFields: nil,
	}, request.SetUserAccessToken("u-GQze1ue1QXC650Z3NYBsga"))
	reqCall.SetBlockId("doxcnIOUiQQCkCSgQK0FF7IKUJh")
	reqCall.SetDocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func listBlocks() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.DocumentBlocks.List(coreCtx, request.SetUserAccessToken("u-G4p3fYOXuJqkwyNOwSDG5g"))
	reqCall.SetDocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf")
	reqCall.SetPageSize(10)

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func listSubBlocks() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.DocumentBlockChildrens.Get(coreCtx, request.SetUserAccessToken("u-GQze1ue1QXC650Z3NYBsga"))
	reqCall.SetDocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf")
	reqCall.SetPageSize(10)
	reqCall.SetBlockId("doxcnIOUiQQCkCSgQK0FF7IKUJh")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func getBlock() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.DocumentBlocks.Get(coreCtx, request.SetUserAccessToken("u-GQze1ue1QXC650Z3NYBsga"))
	reqCall.SetDocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf")
	reqCall.SetBlockId("doxcnIOUiQQCkCSgQK0FF7IKUJh")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func batchDelBlock() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := docxService.DocumentBlockChildrens.BatchDelete(coreCtx, &docx.DocumentBlockChildrenBatchDeleteReqBody{
		StartIndex:      0,
		EndIndex:        1,
		ForceSendFields: []string{"StartIndex"},
	}, request.SetUserAccessToken("u-GQze1ue1QXC650Z3NYBsga"))
	reqCall.SetDocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf")
	reqCall.SetBlockId("doxcnIOUiQQCkCSgQK0FF7IKUJh")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}
