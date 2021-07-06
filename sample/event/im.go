package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	eventhttp "github.com/larksuite/oapi-sdk-go/event/http"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	im "github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func main() {

	// for redis store and logrus
	// var conf = configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
	// var conf = configs.TestConfig("https://open.feishu.cn")
	var conf = configs.TestConfig(core.DomainFeiShu)

	im.SetMessageReceiveEventHandler(conf, func(ctx *core.Context, event *im.MessageReceiveEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(event))
		return nil
	})

	g := gin.Default()
	g.POST("/webhook/event", func(context *gin.Context) {
		eventhttp.Handle(conf, context.Request, context.Writer)
	})
	err := g.Run(":8089")
	if err != nil {
		fmt.Println(err)
	}

}
