package configs

import (
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
)

// for Cutome APP（企业自建应用）
var appSettings = config.GetInternalAppSettingsByEnv()

func TestConfigWithLogrusAndRedisStore(domain constants.Domain) *config.Config {
	logger := Logrus{}
	store := NewRedisStore()
	return config.NewConfig(domain, appSettings, logger, log.LevelDebug, store)
}

func TestConfig(domain constants.Domain) *config.Config {
	return config.NewConfigWithDefaultStore(domain, appSettings, log.NewDefaultLogger(), log.LevelDebug)
}
