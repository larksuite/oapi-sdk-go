package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	speech_to_text "github.com/larksuite/oapi-sdk-go/service/speech_to_text/v1"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var speechToTextService = speech_to_text.NewService(configs.TestConfig(core.DomainFeiShu))

func main() {
	testSpeechFileRecognize()
}

func testSpeechFileRecognize() {
	coreCtx := core.WrapContext(context.Background())
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
		e := err.(*response.Error)
		fmt.Printf(tools.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", tools.Prettify(result))
}
