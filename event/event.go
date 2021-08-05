package event

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/core"
	app "github.com/larksuite/oapi-sdk-go/event/app/v1"
	"github.com/larksuite/oapi-sdk-go/event/core/handlers"
	"github.com/larksuite/oapi-sdk-go/event/core/model"

	"sync"
)

var once sync.Once

func SetTypeHandler(conf core.Config, eventType string, handler handlers.Handler) {
	handlers.SetTypeHandler(conf, eventType, handler)
}

// Deprecated, please use `SetTypeCallback`
func SetTypeHandler2(conf core.Config, eventType string, callback func(ctx *core.Context, event map[string]interface{}) error) {
	SetTypeHandler(conf, eventType, &defaultHandler{callback: callback})
}

func SetTypeCallback(conf core.Config, eventType string, callback func(ctx *core.Context, event map[string]interface{}) error) {
	SetTypeHandler(conf, eventType, &defaultHandler{callback: callback})
}

type defaultHandler struct {
	callback func(ctx *core.Context, event map[string]interface{}) error
}

func (h *defaultHandler) GetEvent() interface{} {
	e := make(map[string]interface{})
	return &e
}

func (h *defaultHandler) Handle(ctx *core.Context, event interface{}) error {
	e := event.(*map[string]interface{})
	return h.callback(ctx, *e)
}

func Handle(conf core.Config, request *core.OapiRequest) *core.OapiResponse {
	once.Do(func() {
		app.SetAppTicketEventHandler(conf)
	})
	coreCtx := core.WrapContext(context.TODO())
	coreCtx.Set(core.CtxKeyConfig, conf)
	httpEvent := &model.HTTPEvent{
		Request:  request,
		Response: &core.OapiResponse{},
	}
	handlers.Handle(coreCtx, httpEvent)
	return httpEvent.Response
}
