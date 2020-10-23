package http

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/event/core/handlers"
	"github.com/larksuite/oapi-sdk-go/event/core/model"
	"net/http"
)

func Handle(conf *config.Config, request *http.Request, response http.ResponseWriter) {
	coreCtx := core.WarpContext(request.Context())
	conf.WithContext(coreCtx)
	httpEvent := &model.HTTPEvent{
		HTTPRequest:  request,
		HTTPResponse: response,
	}
	handlers.Handle(coreCtx, httpEvent)
}
