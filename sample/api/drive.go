package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	drivev1 "github.com/larksuite/oapi-sdk-go/service/drive/v1"
	"hash/adler32"
	"io"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(constants.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var driveService = drivev1.NewService(configs.TestConfig(constants.DomainFeiShu))

func main() {
	testFileUploadAll()
	testFileUploadPart()
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

func testFileUploadPart() {
	coreCtx := core.WrapContext(context.Background())
	userAccessTokenOptFn := request.SetUserAccessToken("[user_access_token]")
	fileSize := 1024

	// upload prepare
	uploadPrepareReqCall := driveService.Files.UploadPrepare(coreCtx, &drivev1.UploadInfo{
		FileName:   fmt.Sprintf("[file_name]"),
		ParentType: "explorer",
		ParentNode: "[folder_token]",
		Size:       fileSize,
	}, userAccessTokenOptFn)

	uploadPrepareResult, err := uploadPrepareReqCall.Do()
	fmt.Printf("[upload prepare] request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("[upload prepare] http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Println(tools.Prettify(e))
		return
	}

	fmt.Printf("[upload prepare] reault:%s", tools.Prettify(uploadPrepareResult))

	// upload part
	uploadedBlockNum := 0
	for i := 0; i < uploadPrepareResult.BlockNum; i++ {
		uploadPartReqCall := driveService.Files.UploadPart(coreCtx, userAccessTokenOptFn)
		uploadPartReqCall.SetUploadId(uploadPrepareResult.UploadId)
		uploadPartReqCall.SetSeq(i)
		//uploadPartReqCall.Set
		// 最后一块
		blockSize := uploadPrepareResult.BlockSize
		if i == (uploadPrepareResult.BlockNum - 1) {
			blockSize = fileSize - (i * uploadPrepareResult.BlockNum)
		}
		uploadPartReqCall.SetSize(blockSize)
		fileContent := createRandomFileData(int64(blockSize))
		file := request.NewFile().SetContent(fileContent)

		uploadPartReqCall.SetFile(file)
		uploadPartReqCall.SetChecksum(fmt.Sprintf("%d", adler32.Checksum(fileContent)))

		result, err := uploadPartReqCall.Do()
		fmt.Printf("[upload part[%d]] request_id:%s", i, coreCtx.GetRequestID())
		fmt.Printf("[upload part[%d]] http status code:%d", i, coreCtx.GetHTTPStatusCode())
		if err != nil {
			e := err.(*response.Error)
			fmt.Println(tools.Prettify(e))
			return
		}
		uploadedBlockNum++
		fmt.Printf("[upload part[%d]] reault:%s", i, tools.Prettify(result))
	}

	// upload finish
	uploadFinishReqCall := driveService.Files.UploadFinish(coreCtx, &drivev1.FileUploadFinishReqBody{
		UploadId: uploadPrepareResult.UploadId,
		BlockNum: uploadedBlockNum,
	}, userAccessTokenOptFn)

	uploadFinishResult, err := uploadFinishReqCall.Do()
	fmt.Printf("[upload finish] request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("[upload finish] http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Println(tools.Prettify(e))
		return
	}

	fmt.Printf("[upload finish] reault:%s", tools.Prettify(uploadFinishResult))

	if len(uploadFinishResult.FileToken) == 0 {
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
