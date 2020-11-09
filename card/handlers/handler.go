package handlers

import (
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

type Handler func(*core.Context, *model.Card) (interface{}, error)

var AppID2Handler = make(map[string]Handler)

func getHandler(conf *config.Config) (Handler, bool) {
	h, ok := AppID2Handler[conf.GetAppSettings().AppID]
	return h, ok
}
