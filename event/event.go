package event

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	coremodel "github.com/larksuite/oapi-sdk-go/core/model"
	app "github.com/larksuite/oapi-sdk-go/event/app/v1"
	"github.com/larksuite/oapi-sdk-go/event/core/handlers"
	"github.com/larksuite/oapi-sdk-go/event/core/model"
	"sync"
)

var once sync.Once

func SetTypeHandler(conf *config.Config, eventType string, handler handlers.Handler) {
	handlers.SetTypeHandler(conf, eventType, handler)
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
	once.Do(func() {
		app.SetAppTicketEventHandler(conf)
	})
	coreCtx := core.WarpContext(request.Ctx)
	conf.WithContext(coreCtx)
	httpEvent := &model.HTTPEvent{
		Request:  request,
		Response: &coremodel.OapiResponse{},
	}
	handlers.Handle(coreCtx, httpEvent)
	return httpEvent.Response
}
