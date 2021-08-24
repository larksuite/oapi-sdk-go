package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
)

func main() {
	// for redis store and logrus
	// var conf = configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
	// var conf = configs.TestConfig("https://open.feishu.cn")
	var conf = configs.TestConfig(core.DomainFeiShu)

	application.SetAppOpenEventHandler(conf, func(coreCtx *core.Context, appOpenEvent *application.AppOpenEvent) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(appOpenEvent)
		fmt.Println(tools.Prettify(appOpenEvent))
		return nil
	})

	application.SetAppStatusChangeEventHandler(conf, func(coreCtx *core.Context, appStatusChangeEvent *application.AppStatusChangeEvent) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(appStatusChangeEvent.Event.AppId)
		fmt.Println(appStatusChangeEvent.Event.Status)
		fmt.Println(tools.Prettify(appStatusChangeEvent))
		return nil
	})

	header := make(map[string][]string)
	// from http request header
	req := &core.OapiRequest{
		Ctx:    context.Background(),
		Header: core.NewOapiHeader(header),
		Body:   "", // from http request body
	}
	resp := event.Handle(conf, req)
	fmt.Println(tools.Prettify(resp))
}
