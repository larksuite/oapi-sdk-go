package lark

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/store"
)

var (
	SetAppCredentials      = core.SetAppCredentials
	SetAppEventKey         = core.SetAppEventKey
	SetHelpDeskCredentials = core.SetHelpDeskCredentials
)

type Config = core.Config
type Domain = core.Domain

const (
	DomainFeiShu    Domain = core.DomainFeiShu
	DomainLarkSuite Domain = core.DomainLarkSuite
)

type LogLevel int

const (
	LogLevelDebug LogLevel = LogLevel(core.LoggerLevelDebug)
	LogLevelInfo  LogLevel = LogLevel(core.LoggerLevelInfo)
	LogLevelWarn  LogLevel = LogLevel(core.LoggerLevelWarn)
	LogLevelError LogLevel = LogLevel(core.LoggerLevelError)
)

type AppConfig struct {
	*config.Config
}

func (ac *AppConfig) SetLogger(log log.Logger) {
	ac.Config.GetLogger().SetLogger(log)
}

func (ac *AppConfig) SetLogLevel(level LogLevel) {
	ac.Config.GetLogger().SetLogLevel(log.Level(level))
}

func (ac *AppConfig) SetStore(store store.Store) {
	ac.Config.SetStore(store)
}

func NewInternalAppConfigByEnv(domain Domain) *AppConfig {
	return &AppConfig{core.NewConfig(domain, core.GetInternalAppSettingsByEnv())}
}

func NewISVAppConfigByEnv(domain Domain) *AppConfig {
	return &AppConfig{core.NewConfig(domain, core.GetISVAppSettingsByEnv())}
}

func NewInternalAppConfig(domain Domain, opts ...config.AppSettingsOpt) *AppConfig {
	return &AppConfig{core.NewConfig(domain, core.NewInternalAppSettings(opts...))}
}

func NewISVAppConfig(domain Domain, opts ...config.AppSettingsOpt) *AppConfig {
	return &AppConfig{core.NewConfig(domain, core.NewISVAppSettings(opts...))}
}
