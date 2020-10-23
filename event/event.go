package event

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

type Handler interface {
	GetEvent() interface{}
	Handle(*core.Context, interface{}) error
}

var appID2Type2EventHandler = map[string]map[string]Handler{}

func SetTypeHandler(conf *config.Config, eventType string, handler Handler) {
	appID := conf.GetAppSettings().AppID
	type2EventHandler, ok := appID2Type2EventHandler[appID]
	if !ok {
		type2EventHandler = map[string]Handler{}
		appID2Type2EventHandler[appID] = type2EventHandler
	}
	type2EventHandler[eventType] = handler
}

func SetTypeHandler2(conf *config.Config, eventType string, fn func(ctx *core.Context, event map[string]interface{}) error) {
	SetTypeHandler(conf, eventType, &defaultHandler{fn: fn})
}

func GetType2EventHandler(conf *config.Config) (map[string]Handler, bool) {
	type2EventHandler, ok := appID2Type2EventHandler[conf.GetAppSettings().AppID]
	return type2EventHandler, ok
}

type defaultHandler struct {
	fn func(ctx *core.Context, event map[string]interface{}) error
}

func (h *defaultHandler) GetEvent() interface{} {
	e := make(map[string]interface{})
	return &e
}

func (h *defaultHandler) Handle(ctx *core.Context, event interface{}) error {
	e := event.(*map[string]interface{})
	return h.fn(ctx, *e)
}
