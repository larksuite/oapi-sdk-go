package handlers

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

type Handler interface {
	GetEvent() interface{}
	Handle(*core.Context, interface{}) error
}

var AppID2Type2EventHandler = map[string]map[string]Handler{}

func getType2EventHandler(conf *config.Config) (map[string]Handler, bool) {
	type2EventHandler, ok := AppID2Type2EventHandler[conf.GetAppSettings().AppID]
	return type2EventHandler, ok
}
