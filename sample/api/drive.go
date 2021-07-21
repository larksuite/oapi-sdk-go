package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	drivev1 "github.com/larksuite/oapi-sdk-go/service/drive/v1"
	"hash/adler32"
	"io"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var driveService = drivev1.NewService(configs.TestConfig(core.DomainFeiShu))

func main() {
	testFileUploadAll()
	testMediaBatchGetTmpDownloadURLs()
	testFileDownload()
}
func createRandomFileData(size int64) []byte {
	randomData := make([]byte, size)
	io.ReadFull(rand.Reader, randomData)
	return randomData
}

func testFileUploadAll() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := driveService.Files.UploadAll(coreCtx, request.SetUserAccessToken("[user_access_token]"))

	reqCall.SetParentType("explorer")
	reqCall.SetParentNode("[folder_token]")
	reqCall.SetFileName(fmt.Sprintf("[file_name]"))
	reqCall.SetSize(1024)

	fileContent := createRandomFileData(1024)
	reqCall.SetChecksum(fmt.Sprintf("%d", adler32.Checksum(fileContent)))
	file := request.NewFile()
	file.SetContent(fileContent)
	reqCall.SetFile(file)

	result, err := reqCall.Do()
	fmt.Printf("request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Println(tools.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", tools.Prettify(result))

	if len(result.FileToken) == 0 {
		fmt.Printf("file token is empty")
		return
	}

}

func testMediaBatchGetTmpDownloadURLs() {

	coreCtx := core.WrapContext(context.Background())
	userAccessTokenOptFn := request.SetUserAccessToken("[user_access_token]")

	reqCall := driveService.Medias.BatchGetTmpDownloadUrl(coreCtx, userAccessTokenOptFn)
	reqCall.SetFileTokens([]string{"[file_token]"}...)

	result, err := reqCall.Do()
	fmt.Printf("request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Println(tools.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", tools.Prettify(result))

	if len(result.TmpDownloadUrls) == 0 {
		fmt.Printf("TmpDownloadUrls len invalid")
		return
	}
}

func testFileDownload() {
	coreCtx := core.WrapContext(context.Background())

	reqCall := driveService.Files.Download(coreCtx, request.SetUserAccessToken("[user_access_token]"))

	reqCall.SetFileToken("[file_token]")

	fileContent := bytes.NewBuffer(nil)
	reqCall.SetResponseStream(fileContent)
	_, err := reqCall.Do()
	fmt.Printf("request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Printf(tools.Prettify(e))
		return
	}
}
