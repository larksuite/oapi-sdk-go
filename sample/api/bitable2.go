/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	larkext "github.com/larksuite/oapi-sdk-go/v3/service/ext"
)

func batchAdd(client *lark.Client) {
	client.Bitable.AppTableRecord.BatchCreate(context.Background(), larkbitable.NewBatchCreateAppTableRecordReqBuilder().
		TableId("id").
		Body(larkbitable.NewBatchCreateAppTableRecordReqBodyBuilder().
			Records([]*larkbitable.AppTableRecord{larkbitable.
				NewAppTableRecordBuilder().
				RecordId("").
				Fields(map[string]interface{}{"a": []*larkbitable.Person{larkbitable.NewPersonBuilder().Name("name").Build()}}).
				Build()}).
			Build()).
		Build())
}

func createBitableFile(client *lark.Client) {
	resp, err := client.Ext.DriveExplorer.CreateFile(context.Background(), larkext.NewCreateFileReqBuilder().
		FolderToken("fldcniHf40Vcv1DoEc8SXeuA0Zd").
		Body(larkext.NewCreateFileReqBodyBuilder().
			Title("title").
			Type(larkext.FileTypeBitable).
			Build()).
		Build(), larkcore.WithUserAccessToken("u-1Kg48B3nh96VzeLBgRanoskhlmB1l54biMG010qyw7rm"))

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp.Data))
	fmt.Println(resp.RequestId())
}

func listFileByIterator(client *lark.Client) {
	iter, err := client.Drive.File.ListByIterator(context.Background(), larkdrive.NewListFileReqBuilder().PageSize(5).FolderToken("fldcniHf40Vcv1DoEc8SXeuA0Zd").Build(),
		larkcore.WithUserAccessToken("u-2NJonELO1d.WwANxJ3Q1fL5k2caRglcFr200g5ww22PK"))
	if err != nil {
		fmt.Println(err)
		return
	}
	count := 0
	for {
		hasNext, file, err := iter.Next()
		if err != nil {
			fmt.Println(err)
			break
		}

		if !hasNext {
			break
		}

		fmt.Println(larkcore.Prettify(file))
		time.Sleep(time.Millisecond * 300)
		count++
	}
	fmt.Println(fmt.Sprintf("total :%d", count))
}

func listFile(client *lark.Client) {
	resp, err := client.Drive.File.List(context.Background(), larkdrive.NewListFileReqBuilder().PageSize(100).FolderToken("fldcniHf40Vcv1DoEc8SXeuA0Zd").Build(),
		larkcore.WithUserAccessToken("u-2NJonELO1d.WwANxJ3Q1fL5k2caRglcFr200g5ww22PK"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("total :%d", len(resp.Data.Files)))
}

func TestAppRecordStruct() {
	tableRecord := larkbitable.AppTableRecord{}
	fields := map[string]interface{}{}
	fields["str"] = "string"
	fields["bool"] = false
	fields["listurl1"] = []larkbitable.Url{*larkbitable.NewUrlBuilder().Text("t1").Link("www.baiducom").Build(), *larkbitable.NewUrlBuilder().Text("t2").Link("www.google").Build()}
	fields["liststr"] = []string{"str1", "str2"}
	fields["listperson"] = []larkbitable.Person{*larkbitable.NewPersonBuilder().Name("n1").Id("id1").Email("e1").Build(), *larkbitable.NewPersonBuilder().Name("n2").Id("id2").Email("e2").Build()}
	fields["listattachment"] = []larkbitable.Attachment{*larkbitable.NewAttachmentBuilder().Name("n1").Url("u1").Build(), *larkbitable.NewAttachmentBuilder().Name("n2").Url("url2").Build()}
	tableRecord.Fields = fields

	fmt.Println(tableRecord.BoolField("bool"))
	fmt.Println(larkcore.Prettify(tableRecord.StringField("str")))
	fmt.Println(larkcore.Prettify(tableRecord.ListUrlField("listurl")))
	fmt.Println(larkcore.Prettify(tableRecord.ListStringField("liststr")))
	fmt.Println(larkcore.Prettify(tableRecord.ListPersonField("listperson")))
	fmt.Println(larkcore.Prettify(tableRecord.ListAttachmentField("listattachment")))

	boolField := tableRecord.BoolField("bool")
	strField := tableRecord.StringField("str")
	listUrlField := tableRecord.ListUrlField("listurl")
	listStrField := tableRecord.ListStringField("liststr")
	listPersonField := tableRecord.ListPersonField("listperson")
	listAttachField := tableRecord.ListAttachmentField("listattachment")
	fmt.Println(boolField, strField, listUrlField, listStrField, listPersonField, listAttachField)
}
func main() {
	TestAppRecordStruct()
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug), lark.WithLogReqAtDebug(true))
	listFileByIterator(client)

}
