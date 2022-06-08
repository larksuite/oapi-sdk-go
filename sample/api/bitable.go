package main

import (
	"context"
	"fmt"

	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	bitable "github.com/larksuite/oapi-sdk-go/service/bitable/v1"
	drivev2 "github.com/larksuite/oapi-sdk-go/service/drive_explorer/v2"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var bittableService = bitable.NewService(configs.TestConfig(core.DomainFeiShu))
var driveExploreService = drivev2.NewService(configs.TestConfig(core.DomainFeiShu))

func createFile() {
	coreCtx := core.WrapContext(context.Background())

	reqCall := driveExploreService.Files.Create(coreCtx, &drivev2.FileCreateReqBody{
		Title: "title",
		Type:  "bitable",
	}, request.SetUserAccessToken("u-G4p3fYOXuJqkwyNOwSDG5g"))
	reqCall.SetFolderToken("fldcniHf40Vcv1DoEc8SXeuA0Zd")
	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}
func main() {
	//createFile()
	//getBitableRawContent()
	addField()
}

func tableList() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := bittableService.AppTables.List(coreCtx, request.SetUserAccessToken("u-G4p3fYOXuJqkwyNOwSDG5g"))
	reqCall.SetPageSize(10)
	reqCall.SetAppToken("bascnpApISZqKuO0uEbcOzGcH6b")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func getBitableRawContent() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := bittableService.Apps.Get(coreCtx, request.SetUserAccessToken("u-G4p3fYOXuJqkwyNOwSDG5g"))
	reqCall.SetAppToken("bascnpApISZqKuO0uEbcOzGcH6b")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}

func addField() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := bittableService.AppTableFields.Create(coreCtx, &bitable.AppTableField{
		FieldName: "日期",
		Type:      5,
		Property: &bitable.AppTableFieldProperty{
			DateFormatter: "yyyy/MM/dd HH:mm",
			AutoFill:      false,
		},
	}, request.SetUserAccessToken("u-2rfzaUq7NcN9rAWQ3P4x.r4l7U0w0hC3iww0k1.02aAg"))
	reqCall.SetAppToken("bascnpApISZqKuO0uEbcOzGcH6b")
	reqCall.SetTableId("tbl9a6pWtNVFsSrQ")

	resp, err := reqCall.Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tools.Prettify(resp))
	}
}
