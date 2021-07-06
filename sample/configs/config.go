package configs

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

// for Cutome APP（企业自建应用）
var appSettings = core.GetInternalAppSettingsByEnv()

func TestConfigWithLogrusAndRedisStore(domain core.Domain) *config.Config {
	logger := Logrus{}
	store := NewRedisStore()
	return core.NewConfig(domain, appSettings, core.SetLogger(logger), core.SetLoggerLevel(core.LoggerLevelDebug), core.SetStore(store))
}

func TestConfig(domain core.Domain) *config.Config {
	return core.NewConfig(domain, appSettings, core.SetLoggerLevel(core.LoggerLevelDebug))
}
