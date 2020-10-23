package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/card"
	cardginserver "github.com/larksuite/oapi-sdk-go/card/http/gin"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"path"
)

func main() {

	conf := test.GetOnlineInternalConf()

	card.SetHandler(conf, func(coreCtx *core.Context, card *model.Card) (interface{}, error) {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(tools.Prettify(card.Action))
		return nil, nil
	})

	g := gin.Default()
	cardginserver.Register(path.Join("/", conf.GetAppSettings().AppID, "webhook/card"), conf, g)
	err := g.Run(":8089")
	if err != nil {
		panic(err)
	}
}
