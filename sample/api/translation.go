package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/sample"
	translation "github.com/larksuite/oapi-sdk-go/service/translation/v1"
)

// for redis store and logrus
// sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
// sample.TestConfig("https://open.feishu.cn")
var translationService = translation.NewService(sample.TestConfig(lark.DomainFeiShu))

func main() {
	testTextDetect()
}

func testTextDetect() {
	coreCtx := lark.WrapContext(context.Background())
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
		e := err.(lark.APIError)
		fmt.Printf(lark.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", lark.Prettify(result))
}
