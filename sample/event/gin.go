package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
)

var appConf *lark.AppConfig

func main() {
	appConf = lark.NewInternalAppConfigByEnv(lark.DomainFeiShu)
	appConf.SetLogLevel(lark.LogLevelDebug)

	application.SetAppOpenEventHandler(appConf, func(ctx *lark.Context, appOpenEvent *application.AppOpenEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(appOpenEvent)
		fmt.Println(lark.Prettify(appOpenEvent))
		return nil
	})

	lark.WebHook.SetEventHandler(appConf, "app_status_change", func(ctx *lark.Context, event map[string]interface{}) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(event))
		data := event["event"].(map[string]interface{})
		fmt.Println(lark.Prettify(data))
		return nil
	})

	application.SetAppUninstalledEventHandler(appConf, func(ctx *lark.Context, appUninstalledEvent *application.AppUninstalledEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(appUninstalledEvent))
		return nil
	})

	g := gin.Default()

	g.POST("/webhook/event", func(context *gin.Context) {
		lark.WebHook.EventWebServeHandler(appConf, context.Request, context.Writer)
	})
	err := g.Run(":8089")
	if err != nil {
		panic(err)
	}
}
