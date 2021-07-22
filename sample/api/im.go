package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	im "github.com/larksuite/oapi-sdk-go/service/im/v1"
	"io/ioutil"
	"os"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var imService = im.NewService(configs.TestConfig(core.DomainFeiShu))

func main() {
	testMessageCreate()
	testMessageReactionAll()
	//testFileCreate()
	testFileRead()
}

func testMessageCreate() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := imService.Messages.Create(coreCtx, &im.MessageCreateReqBody{
		// ReceiveId: "b1g6b445",
		ReceiveId: "ou_a11d2bcc7d852afbcaf37e5b3ad01f7e",
		Content:   "{\"text\":\"<at user_id=\\\"ou_a11d2bcc7d852afbcaf37e5b3ad01f7e\\\">Tom</at> test content\"}",
		MsgType:   "text",
	})
	reqCall.SetReceiveIdType("open_id")
	message, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(message))
}

func testMessageReactionAll() {
	coreCtx := core.WrapContext(context.Background())
	createReqCall := imService.MessageReactions.Create(coreCtx, &im.MessageReactionCreateReqBody{
		ReactionType:    &im.Emoji{
			EmojiType:       "JIAYI",
		},
	})
	createReqCall.SetMessageId("om_a8f2294b037bfbd808c9a1a38afaac9d")
	createResp, createErr := createReqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if createErr != nil {
		fmt.Println(tools.Prettify(createErr))
		e := createErr.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	reactionID:= createResp.ReactionId
	fmt.Println(tools.Prettify(createResp))

	listReqCall :=imService.MessageReactions.List(coreCtx)
	listReqCall.SetMessageId("om_a8f2294b037bfbd808c9a1a38afaac9d")
	listReqCall.SetPageSize(1)
	listReqCall.SetReactionType("SMILE")
	listReqCall.SetPageToken("YhljsPiGfUgnVAg9urvRFU6IjGHbxJ2BLswa0r0bYxy7k3ZwdDpAuMABmfMa4aQY")
	listResp, listErr := listReqCall.Do()
	if listErr != nil {
		fmt.Println(tools.Prettify(listErr))
		e := listErr.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(listResp))

	deleteReqCall:=imService.MessageReactions.Delete(coreCtx)
	deleteReqCall.SetMessageId("om_a8f2294b037bfbd808c9a1a38afaac9d")
	deleteReqCall.SetReactionId(reactionID)
	deleteResp, deleteErr :=deleteReqCall.Do()
	if deleteErr != nil {
		fmt.Println(tools.Prettify(deleteErr))
		e := deleteErr.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(deleteResp))
}

func testFileRead() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := imService.Files.Get(coreCtx)
	buf := &bytes.Buffer{}
	reqCall.SetResponseStream(buf)
	reqCall.SetFileKey("file_ec24f8ad-89ea-4bb5-a7e4-c5db35d2925g")
	_, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	fmt.Println(coreCtx.GetHeader())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	err = ioutil.WriteFile("test_download.pdf", buf.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func testFileCreate() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := imService.Files.Create(coreCtx)
	f, err := os.Open("test.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	file := request.NewFile().SetContentStream(f)
	// request.NewFile().SetContent([]byte)
	reqCall.SetFile(file)
	reqCall.SetFileName("test-测试.pdf")
	reqCall.SetFileType("pdf")
	message, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(message))
}
