package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/sample"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
	"net/http"
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

	lark.WebHook.SetEventHandler(conf, "user.created_v2", func(coreCtx *core.Context, event map[string]interface{}) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(lark.Prettify(event))
		return nil
	})

	lark.WebHook.EventWebServeRouter("/webhook/event", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}

}
