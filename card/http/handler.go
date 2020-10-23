package http

import (
	"github.com/larksuite/oapi-sdk-go/card/handlers"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	nethttp "net/http"
)

func Handle(conf *config.Config, request *nethttp.Request, response nethttp.ResponseWriter) {
	coreCtx := core.WarpContext(request.Context())
	conf.WithContext(coreCtx)
	httpCard := &model.HTTPCard{
		HTTPRequest:  request,
		HTTPResponse: response,
	}
	handlers.Handle(coreCtx, httpCard)
}
