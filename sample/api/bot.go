package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	bot "github.com/larksuite/oapi-sdk-go/service/bot/v3"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var botService = bot.NewService(configs.TestConfig(lark.DomainFeiShu))

func main() {
	testBotGet()
}

func testBotGet() {
	coreCtx := lark.WrapContext(context.Background())
	reqCall := botService.Bots.Get(coreCtx)
	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(result.Bot))
}
