package lark

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	. "github.com/larksuite/oapi-sdk-go/event/http"
	"github.com/larksuite/oapi-sdk-go/event/http/native"
	"net/http"
)

func (wb webHook) EventRequestHandle(conf core.Config, request *HTTPRequest) *HTTPResponse {
	resp := event.Handle(conf, &core.OapiRequest{
		Header: core.NewOapiHeader(request.Header),
		Body:   request.Body,
	})
	return resp
}

func (wb webHook) EventWebServeRouter(path string, conf core.Config) {
	native.Register(path, conf)
}

func (wb webHook) EventWebServeHandler(conf core.Config, request *http.Request, response http.ResponseWriter) {
	Handle(conf, request, response)
}

func (wb webHook) SetEventHandler(conf core.Config, eventType string, handler func(ctx *Context, event map[string]interface{}) error) {
	event.SetTypeCallback(conf, eventType, func(ctx *core.Context, event map[string]interface{}) error {
		return handler(ctx, event)
	})
}
