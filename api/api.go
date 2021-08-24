package api

import (
	"github.com/larksuite/oapi-sdk-go/api/core/handlers"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

func Send(ctx *core.Context, conf *config.Config, req *request.Request) error {
	ctx.Set(config.CtxKeyConfig, conf)
	req.WithContext(ctx)
	handlers.Handle(ctx, req)
	if req.Err == nil {
		return nil
	}
	return response.ToError(req.Err)
}
