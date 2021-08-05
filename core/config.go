package core

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/store"
)

type Config interface {
	GetDomain() string
	GetAppSettings() *config.AppSettings
	GetLogger() *log.LoggerProxy
	GetStore() store.Store
	GetHelpDeskAuthorization() string
}

var CtxKeyConfig = "-----CtxKeyConfig"

func GetConfigByCtx(ctx context.Context) Config {
	c := ctx.Value(CtxKeyConfig)
	return c.(Config)
}
