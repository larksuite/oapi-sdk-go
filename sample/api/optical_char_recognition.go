package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	optical_char_recognition "github.com/larksuite/oapi-sdk-go/service/optical_char_recognition/v1"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var opticalCharRecognitionService = optical_char_recognition.NewService(configs.TestConfig(core.DomainFeiShu))

func main() {
	testImageBasicRecognize()
}

func testImageBasicRecognize() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := opticalCharRecognitionService.Images.BasicRecognize(coreCtx, &optical_char_recognition.ImageBasicRecognizeReqBody{
		Image: "base64 image",
	})
	result, err := reqCall.Do()
	fmt.Printf("request_id:%s\n", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Printf(tools.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", tools.Prettify(result))
}
