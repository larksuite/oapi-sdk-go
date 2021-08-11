package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/sample"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
)

func main() {
	// for redis store and logrus
	// var conf = sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
	// var conf = sample.TestConfig("https://open.feishu.cn")
	var conf = sample.TestConfig(lark.DomainFeiShu)

	application.SetAppOpenEventHandler(conf, func(coreCtx *core.Context, appOpenEvent *application.AppOpenEvent) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(appOpenEvent)
		fmt.Println(lark.Prettify(appOpenEvent))
		return nil
	})

	application.SetAppStatusChangeEventHandler(conf, func(coreCtx *core.Context, appStatusChangeEvent *application.AppStatusChangeEvent) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(appStatusChangeEvent.Event.AppId)
		fmt.Println(appStatusChangeEvent.Event.Status)
		fmt.Println(lark.Prettify(appStatusChangeEvent))
		return nil
	})

	header := make(map[string][]string)
	resp := lark.WebHook.EventRequestHandle(conf, &lark.HTTPRequest{
		Header: header,
		Body:   "", // from http request body
	})
	fmt.Println(lark.Prettify(resp))
}
