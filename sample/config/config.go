package config

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
)

func GetConfig(domain constants.Domain, appSettings *config.AppSettings, level log.Level) *config.Config {
	logger := Logrus{}
	store := NewRedisStore()
	coreCtx := core.WrapContext(context.Background())
	coreCtx.GetHTTPStatusCode()
	return config.NewConfig(constants.DomainLarkSuite, appSettings, logger, level, store)
}
