package lark

import (
	"github.com/larksuite/oapi-sdk-go/card"
	. "github.com/larksuite/oapi-sdk-go/card/http"
	"github.com/larksuite/oapi-sdk-go/card/http/native"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"net/http"
)

type CardAction = model.Card

func (wb webHook) CardRequestHandle(conf core.Config, request *HTTPRequest) *HTTPResponse {
	resp := card.Handle(conf, &core.OapiRequest{
		Header: core.NewOapiHeader(request.Header),
		Body:   request.Body,
	})
	return resp
}

func (wb webHook) CardWebServeRouter(path string, conf core.Config) {
	native.Register(path, conf)
}

func (wb webHook) CardWebServeHandler(conf core.Config, request *http.Request, response http.ResponseWriter) {
	Handle(conf, request, response)
}

func (wb webHook) SetCardActionHandler(conf core.Config, handler func(ctx *Context, action *CardAction) (interface{}, error)) {
	card.SetHandler(conf, handler)
}
