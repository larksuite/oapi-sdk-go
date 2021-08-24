package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	eventhttpserver "github.com/larksuite/oapi-sdk-go/event/http/native"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
	"net/http"
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

	eventhttpserver.Register("/webhook/event", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}

}
