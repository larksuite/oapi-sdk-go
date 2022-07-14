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
	}

	// 业务数据处理
	fmt.Println(larkcore.Prettify(resp.Data))
}

func listBlocks(client *lark.Client) {
	resp, err := client.Docx.DocumentBlock.List(context.Background(),
		larkdocx.NewListDocumentBlockReqBuilder().
			DocumentId("doxcnku1W0IhiZBDPkxlEVSn6Tf").
			PageSize(100).
			Build(), larkcore.WithUserAccessToken("u-3vEh2SpiF2WoJzYJOdiGKQ41mJrQ1hebh0G0hg.02CgW"),
	)

	if err != nil {
		fmt.Println(err)
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
		larkcore.WithUserAccessToken("u-11ETll3Kd1O8NxVwd_uVVN0hnoUAlhcbWi00kg.yyIsw"))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.RequestId())
	fmt.Println(larkcore.Prettify(resp))

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
			Build(), larkcore.WithUserAccessToken("u-11ETll3Kd1O8NxVwd_uVVN0hnoUAlhcbWi00kg.yyIsw"),
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
