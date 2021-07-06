package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	translation "github.com/larksuite/oapi-sdk-go/service/translation/v1"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var translationService = translation.NewService(configs.TestConfig(core.DomainFeiShu))

func main() {
	testTextDetect()
}

func testTextDetect() {
	coreCtx := core.WrapContext(context.Background())
	reqCall := translationService.Texts.Translate(coreCtx, &translation.TextTranslateReqBody{
		SourceLanguage: "zh",
		Text:           "测试",
		TargetLanguage: "en",
		Glossary: []*translation.Term{
			{
				From: "测",
				To:   "test",
			},
		},
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
