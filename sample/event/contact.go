package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/sample"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
)

func main() {

	// for redis store and logrus
	// var conf = sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
	// var conf = sample.TestConfig("https://open.feishu.cn")
	var conf = sample.TestConfig(lark.DomainFeiShu)

	contact.SetDepartmentCreatedEventHandler(conf, func(ctx *core.Context, event *contact.DepartmentCreatedEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(event))
		return nil
	})

	contact.SetUserCreatedEventHandler(conf, func(ctx *core.Context, event *contact.UserCreatedEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(event))
		return nil
	})

	contact.SetUserUpdatedEventHandler(conf, func(ctx *core.Context, event *contact.UserUpdatedEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(event))
		return nil
	})

	g := gin.Default()
	g.POST("/webhook/event", func(context *gin.Context) {
		lark.WebHook.EventWebServeHandler(conf, context.Request, context.Writer)
	})
	err := g.Run(":8089")
	if err != nil {
		fmt.Println(err)
	}

}
