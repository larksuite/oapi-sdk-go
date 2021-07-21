package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/card"
	cardhttp "github.com/larksuite/oapi-sdk-go/card/http"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
)

func main() {

	// for redis store and logrus
	// var conf = configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
	// var conf = configs.TestConfig("https://open.feishu.cn")
	var conf = configs.TestConfig(core.DomainFeiShu)

	card.SetHandler(conf, func(coreCtx *core.Context, card *model.Card) (interface{}, error) {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(tools.Prettify(card.Action))
		return nil, nil
	})

	g := gin.Default()
	g.POST("/webhook/card", func(context *gin.Context) {
		cardhttp.Handle(conf, context.Request, context.Writer)
	})
	err := g.Run(":8089")
	if err != nil {
		panic(err)
	}
}
