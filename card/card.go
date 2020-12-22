package card

import (
	"github.com/larksuite/oapi-sdk-go/card/handlers"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	coremodel "github.com/larksuite/oapi-sdk-go/core/model"
)

func SetHandler(conf *config.Config, handler handlers.Handler) {
	handlers.AppID2Handler[conf.GetAppSettings().AppID] = handler
}

func Handle(conf *config.Config, request *coremodel.OapiRequest) *coremodel.OapiResponse {
	coreCtx := core.WrapContext(request.Ctx)
	conf.WithContext(coreCtx)
	httpCard := &model.HTTPCard{
		Request:  request,
		Response: &coremodel.OapiResponse{},
	}
	handlers.Handle(coreCtx, httpCard)
	return httpCard.Response
}
