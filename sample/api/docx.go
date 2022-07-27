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
	"net/http"
	"os"
	"time"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	"github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func createDocument(client *lark.Client) {
	// 自定义请求headers
	header := make(http.Header)
	header.Add("k1", "v1")
	header.Add("k2", "v2")

	// 发起请求
	resp, err := client.Docx.Document.Create(context.Background(), larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			FolderToken("token").
			Title("title").
			Build(),
		).
		Build(),
		larkcore.WithUserAccessToken("userToken"), // 设置用户Token
		larkcore.WithHeaders(header),              // 设置自定义headers
	)

	//处理错误
	if err != nil {
		// 处理err
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	// 业务数据处理
	fmt.Println(larkcore.Prettify(resp.Data))
}

func listBlocks(client *lark.Client) {
	resp, err := client.Docx.DocumentBlock.List(context.Background(),
		larkdocx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(100).
			Build(), larkcore.WithUserAccessToken("u-1C.E95YFlf2HqXDz4kcNjx5lhNtMh5CxqMG0l0a00yWy"),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(resp.RequestId())
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(len(resp.Data.Items))

}

func downloadFile(client *lark.Client) {
	resp, err := client.Drive.File.Download(context.Background(),
		larkdrive.NewDownloadFileReqBuilder().
			FileToken("boxcnTrRml0GB9E3NFDEyNtMeOb").
			Build(),
		larkcore.WithUserAccessToken("u-1C.E95YFlf2HqXDz4kcNjx5lhNtMh5CxqMG0l0a00yWy"))

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func listBlocksIter(client *lark.Client) {
	var count = 0

	defer func() {
		fmt.Println(count)

	}()

	iter, err := client.Docx.DocumentBlock.ListByIterator(context.Background(),
		larkdocx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(1).
			Limit(3).
			Build(), larkcore.WithUserAccessToken("u-1C.E95YFlf2HqXDz4kcNjx5lhNtMh5CxqMG0l0a00yWy"),
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

		fmt.Println(larkcore.Prettify(block))
		time.Sleep(time.Second)
		count++
	}

}

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")

	feishuClient := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	downloadFile(feishuClient)
	//listBlocks(feishuClient)
	//listBlocksIter(feishuClient)
}
