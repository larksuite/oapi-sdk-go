package core

import (
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/store"
)

var (
	GetISVAppSettingsByEnv      = config.GetISVAppSettingsByEnv
	GetInternalAppSettingsByEnv = config.GetInternalAppSettingsByEnv
	NewISVAppSettings           = config.NewISVAppSettingsByOpts
	NewInternalAppSettings      = config.NewInternalAppSettingsByOpts
	SetAppCredentials           = config.SetAppCredentials
	SetAppEventKey              = config.SetAppEventKey
	SetHelpDeskCredentials      = config.SetHelpDeskCredentials
)

type LoggerLevel int

const (
	LoggerLevelDebug LoggerLevel = LoggerLevel(log.LevelDebug)
	LoggerLevelInfo  LoggerLevel = LoggerLevel(log.LevelInfo)
	LoggerLevelWarn  LoggerLevel = LoggerLevel(log.LevelWarn)
	LoggerLevelError LoggerLevel = LoggerLevel(log.LevelError)
)

type configOpt struct {
	logger   log.Logger
	logLevel log.Level
	store    store.Store
}

type Domain string

const (
	DomainFeiShu    Domain = Domain(constants.DomainFeiShu)
	DomainLarkSuite Domain = Domain(constants.DomainLarkSuite)
)

func NewConfig(domain Domain, appSettings *config.AppSettings, opts ...ConfigOpt) *config.Config {
	configOpt := &configOpt{
		logLevel: log.LevelError,
		logger:   log.NewDefaultLogger(),
	}
	for _, opt := range opts {
		opt(configOpt)
	}
	if configOpt.store == nil {
		return config.NewConfigWithDefaultStore(constants.Domain(domain), appSettings, configOpt.logger, configOpt.logLevel)
	}
	return config.NewConfig(constants.Domain(domain), appSettings, configOpt.logger, configOpt.logLevel, configOpt.store)
}

type ConfigOpt func(o *configOpt)

func SetLogger(logger log.Logger) func(o *configOpt) {
	return func(o *configOpt) {
		o.logger = logger
	}
}

func SetLoggerLevel(logLevel LoggerLevel) func(o *configOpt) {
	return func(o *configOpt) {
		o.logLevel = log.Level(logLevel)
	}
}

func SetStore(store store.Store) func(o *configOpt) {
	return func(o *configOpt) {
		o.store = store
	}
}
