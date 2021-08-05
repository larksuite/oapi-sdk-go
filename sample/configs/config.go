package configs

import (
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
)

func TestConfigWithLogrusAndRedisStore(domain lark.Domain) core.Config {
	conf := lark.NewInternalAppConfigByEnv(domain)
	conf.SetLogger(Logrus{})
	conf.SetLogLevel(lark.LogLevelDebug)
	conf.SetStore(NewRedisStore())
	return conf
}

func TestConfig(domain lark.Domain) core.Config {
	conf := lark.NewInternalAppConfigByEnv(domain)
	conf.SetLogLevel(lark.LogLevelDebug)
	return conf
}
