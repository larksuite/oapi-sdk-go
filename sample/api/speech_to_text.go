package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/sample"
	speech_to_text "github.com/larksuite/oapi-sdk-go/service/speech_to_text/v1"
)

// for redis store and logrus
// sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
// sample.TestConfig("https://open.feishu.cn")
var speechToTextService = speech_to_text.NewService(sample.TestConfig(lark.DomainFeiShu))

func main() {
	testSpeechFileRecognize()
}

func testSpeechFileRecognize() {
	coreCtx := lark.WrapContext(context.Background())
	reqCall := speechToTextService.Speechs.FileRecognize(coreCtx, &speech_to_text.SpeechFileRecognizeReqBody{
		Speech: &speech_to_text.Speech{
			Speech: "base64 后的音频内容",
		},
		Config: &speech_to_text.FileConfig{
			FileId:     "qwe12dd34567890w",
			Format:     "pcm",
			EngineType: "16k_auto",
		},
	})
	result, err := reqCall.Do()
	fmt.Printf("header:%s\n", coreCtx.GetHeader())
	fmt.Printf("request_id:%s\n", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(lark.APIError)
		fmt.Printf(lark.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", lark.Prettify(result))
}
