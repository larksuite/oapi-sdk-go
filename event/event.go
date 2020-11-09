package event

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	coremodel "github.com/larksuite/oapi-sdk-go/core/model"
	"github.com/larksuite/oapi-sdk-go/event/core/handlers"
	"github.com/larksuite/oapi-sdk-go/event/core/model"
)

func SetTypeHandler(conf *config.Config, eventType string, handler handlers.Handler) {
	appID := conf.GetAppSettings().AppID
	type2EventHandler, ok := handlers.AppID2Type2EventHandler[appID]
	if !ok {
		type2EventHandler = map[string]handlers.Handler{}
		handlers.AppID2Type2EventHandler[appID] = type2EventHandler
	}
	type2EventHandler[eventType] = handler
}

func SetTypeHandler2(conf *config.Config, eventType string, fn func(ctx *core.Context, event map[string]interface{}) error) {
	SetTypeHandler(conf, eventType, &defaultHandler{fn: fn})
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

func Handle(conf *config.Config, request *coremodel.OapiRequest) *coremodel.OapiResponse {
	coreCtx := core.WarpContext(request.Ctx)
	conf.WithContext(coreCtx)
	httpEvent := &model.HTTPEvent{
		Request:  request,
		Response: &coremodel.OapiResponse{},
	}
	handlers.Handle(coreCtx, httpEvent)
	return httpEvent.Response
}
