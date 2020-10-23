package card

import (
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

type Handler func(*core.Context, *model.Card) (interface{}, error)

var appID2Handler = make(map[string]Handler)

func SetHandler(conf *config.Config, handler Handler) {
	appID2Handler[conf.GetAppSettings().AppID] = handler
}

func GetHandler(conf *config.Config) (Handler, bool) {
	h, ok := appID2Handler[conf.GetAppSettings().AppID]
	return h, ok
}
