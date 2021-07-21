package config

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/store"
)

var CtxKeyConfig = "-----CtxKeyConfig"

type Config struct {
	domain                constants.Domain
	appSettings           *AppSettings
	store                 store.Store // store
	logger                log.Logger  // logger
	helpDeskAuthorization string
}

func NewTestConfig(domain constants.Domain, appSettings *AppSettings) *Config {
	return NewConfigWithDefaultStore(domain, appSettings, log.NewDefaultLogger(), log.LevelDebug)
}

func NewConfigWithDefaultStore(domain constants.Domain, appSettings *AppSettings, logger log.Logger, logLevel log.Level) *Config {
	loggerProxy := log.NewLoggerProxy(logLevel, logger)
	conf := &Config{
		domain:                domain,
		appSettings:           appSettings,
		store:                 store.NewDefaultStoreWithLog(loggerProxy),
		logger:                loggerProxy,
		helpDeskAuthorization: appSettings.helpDeskAuthorization(),
	}
	return conf
}

func NewConfig(domain constants.Domain, appSettings *AppSettings, logger log.Logger, logLevel log.Level, store store.Store) *Config {
	loggerProxy := log.NewLoggerProxy(logLevel, logger)
	conf := &Config{
		domain:                domain,
		appSettings:           appSettings,
		store:                 store,
		logger:                loggerProxy,
		helpDeskAuthorization: appSettings.helpDeskAuthorization(),
	}
	return conf
}

func (c *Config) GetDomain() string {
	return string(c.domain)
}

func (c *Config) GetAppSettings() *AppSettings {
	return c.appSettings
}

func (c *Config) GetLogger() log.Logger {
	return c.logger
}

func (c *Config) GetStore() store.Store {
	return c.store
}

func (c *Config) GetHelpDeskAuthorization() string {
	return c.helpDeskAuthorization
}

func ByCtx(ctx context.Context) *Config {
	c := ctx.Value(CtxKeyConfig)
	return c.(*Config)
}
