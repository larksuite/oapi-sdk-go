package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/sample"
	optical_char_recognition "github.com/larksuite/oapi-sdk-go/service/optical_char_recognition/v1"
)

// for redis store and logrus
// sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
// sample.TestConfig("https://open.feishu.cn")
var opticalCharRecognitionService = optical_char_recognition.NewService(sample.TestConfig(lark.DomainFeiShu))

func main() {
	testImageBasicRecognize()
}

func testImageBasicRecognize() {
	coreCtx := lark.WrapContext(context.Background())
	reqCall := opticalCharRecognitionService.Images.BasicRecognize(coreCtx, &optical_char_recognition.ImageBasicRecognizeReqBody{
		Image: "base64 image",
	})
	result, err := reqCall.Do()
	fmt.Printf("request_id:%s\n", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(lark.APIError)
		fmt.Printf(lark.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", lark.Prettify(result))
}
