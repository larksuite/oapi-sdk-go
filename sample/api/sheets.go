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

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larksheets "github.com/larksuite/oapi-sdk-go/v3/service/sheets/v3"
)

func createSheet(client *lark.Client) {
	// 用户权限创建表格
	resp, err := client.Sheets.Spreadsheet.Create(context.Background(), larksheets.
		NewCreateSpreadsheetReqBuilder().
		Spreadsheet(larksheets.NewSpreadsheetBuilder().
			FolderToken("fldcniHf40Vcv1DoEc8SXeuA0Zd").
			Title("title").
			Build()).
		Build(),
		larkcore.WithUserAccessToken("u-2NJonELO1d.WwANxJ3Q1fL5k2caRglcFr200g5ww22PK"))

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

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug), lark.WithLogReqAtDebug(true))

	createSheet(client)
}
