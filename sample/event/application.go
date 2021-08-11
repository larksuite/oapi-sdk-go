package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
	"net/http"
)

func main() {

	// for redis store and logrus
	// var conf = sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
	// var conf = sample.TestConfig("https://open.feishu.cn")
	var conf = lark.NewInternalAppConfigByEnv(lark.DomainFeiShu)

	application.SetAppOpenEventHandler(conf, func(ctx *lark.Context, appOpenEvent *application.AppOpenEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(appOpenEvent)
		fmt.Println(lark.Prettify(appOpenEvent))
		return nil
	})

	application.SetAppStatusChangeEventHandler(conf, func(ctx *lark.Context, appStatusChangeEvent *application.AppStatusChangeEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(appStatusChangeEvent.Event.AppId)
		fmt.Println(appStatusChangeEvent.Event.Status)
		fmt.Println(lark.Prettify(appStatusChangeEvent))
		return nil
	})

	lark.WebHook.EventWebServeRouter("/webhook/event", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}

}
