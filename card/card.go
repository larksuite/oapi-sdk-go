package card

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/card/handlers"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
)

func SetHandler(conf core.Config, handler handlers.Handler) {
	handlers.AppID2Handler[conf.GetAppSettings().AppID] = handler
}

func Handle(conf core.Config, request *core.OapiRequest) *core.OapiResponse {
	coreCtx := core.WrapContext(context.TODO())
	coreCtx.Set(core.CtxKeyConfig, conf)
	httpCard := &model.HTTPCard{
		Request:  request,
		Response: &core.OapiResponse{},
	}
	handlers.Handle(coreCtx, httpCard)
	return httpCard.Response
}
