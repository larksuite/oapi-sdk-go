package card

import (
	"github.com/larksuite/oapi-sdk-go/card/handlers"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

func SetHandler(conf *config.Config, handler handlers.Handler) {
	handlers.AppID2Handler[conf.GetAppSettings().AppID] = handler
}

func Handle(conf *config.Config, request *core.OapiRequest) *core.OapiResponse {
	coreCtx := core.WrapContext(request.Ctx)
	coreCtx.Set(config.CtxKeyConfig, conf)
	httpCard := &model.HTTPCard{
		Request:  request,
		Response: &core.OapiResponse{},
	}
	handlers.Handle(coreCtx, httpCard)
	return httpCard.Response
}
