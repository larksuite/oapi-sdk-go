package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/card"
	cardhttpserver "github.com/larksuite/oapi-sdk-go/card/http/native"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"net/http"
	"path"
)

func main() {

	conf := test.GetOnlineInternalConf()

	card.SetHandler(conf, func(ctx *core.Context, card *model.Card) (interface{}, error) {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(card)
		return nil, nil
	})

	cardhttpserver.Register(path.Join("/", conf.GetAppSettings().AppID, "webhook/card"), conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}

}
