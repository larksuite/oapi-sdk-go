package http

import (
	"github.com/larksuite/oapi-sdk-go/card/handlers"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	coremodel "github.com/larksuite/oapi-sdk-go/core/model"
	"net/http"
)

func Handle(conf *config.Config, request *http.Request, response http.ResponseWriter) {
	req, err := coremodel.ToOapiRequest(request)
	if err != nil {
		err = coremodel.NewOapiResponseOfErr(err).WriteTo(response)
		if err != nil {
			conf.GetLogger().Error(req.Ctx, err)
		}
		return
	}
	err = Handle2(conf, req).WriteTo(response)
	if err != nil {
		conf.GetLogger().Error(req.Ctx, err)
	}
}

func Handle2(conf *config.Config, request *coremodel.OapiRequest) *coremodel.OapiResponse {
	coreCtx := core.WarpContext(request.Ctx)
	conf.WithContext(coreCtx)
	httpEvent := &model.HTTPCard{
		Request:  request,
		Response: &coremodel.OapiResponse{},
	}
	handlers.Handle(coreCtx, httpEvent)
	return httpEvent.Response
}
