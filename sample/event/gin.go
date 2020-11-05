package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/event"
	eventginserver "github.com/larksuite/oapi-sdk-go/event/http/gin"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
)

func main() {

	var conf = test.GetISVConf("online")

	application.SetAppOpenEventHandler(conf, func(ctx *core.Context, appOpenEvent *application.AppOpenEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(appOpenEvent)
		fmt.Println(tools.Prettify(appOpenEvent))
		return nil
	})

	/*
		application.SetAppStatusChangeEventHandler(conf, func(ctx *core.Context, appStatusChangeEvent *application.AppStatusChangeEvent) error {
			fmt.Println(ctx.GetRequestID())
			fmt.Println(appStatusChangeEvent.Event.AppId)
			fmt.Println(appStatusChangeEvent.Event.Status)
			fmt.Println(tools.Prettify(appStatusChangeEvent))
			return nil
		})
	*/
	event.SetTypeHandler2(conf, "app_status_change", func(ctx *core.Context, event map[string]interface{}) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(event))
		data := event["event"].(map[string]interface{})
		fmt.Println(tools.Prettify(data))
		return nil
	})

	application.SetAppUninstalledEventHandler(conf, func(ctx *core.Context, appUninstalledEvent *application.AppUninstalledEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(appUninstalledEvent))
		return nil
	})

	application.SetAppUninstalledEventHandler(conf, func(ctx *core.Context, appUninstalledEvent *application.AppUninstalledEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(appUninstalledEvent))
		return nil
	})

	application.SetOrderPaidEventHandler(conf, func(ctx *core.Context, orderPaidEvent *application.OrderPaidEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(orderPaidEvent))
		return nil
	})

	g := gin.Default()
	eventginserver.Register("/webhook/event", conf, g)
	err := g.Run(":8089")
	if err != nil {
		panic(err)
	}
}
