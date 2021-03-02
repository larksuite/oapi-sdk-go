package test

import (
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"os"
	"strings"
)

func domainFeiShu(env string) string {
	return os.Getenv(env + "_FEISHU_DOMAIN")
}

func getISVAppSettings(env string) *config.AppSettings {
	appID, appSecret, verificationToken, encryptKey := os.Getenv(env+"_ISV_APP_ID"),
		os.Getenv(env+"_ISV_APP_SECRET"), os.Getenv(env+"_ISV_VERIFICATION_TOKEN"), os.Getenv(env+"_ISV_ENCRYPT_KEY")
	return config.NewISVAppSettings(appID, appSecret, verificationToken, encryptKey)
}

func getInternalAppSettings(env string) *config.AppSettings {
	appID, appSecret, verificationToken, encryptKey := os.Getenv(env+"_INTERNAL_APP_ID"),
		os.Getenv(env+"_INTERNAL_APP_SECRET"), os.Getenv(env+"_INTERNAL_VERIFICATION_TOKEN"), os.Getenv(env+"_INTERNAL_ENCRYPT_KEY")
	return config.NewInternalAppSettings(appID, appSecret, verificationToken, encryptKey)
}

func GetISVConf(env string) *config.Config {
	env = strings.ToUpper(env)
	return config.NewTestConfig(getDomain(env), getISVAppSettings(env))
}

func GetInternalConf(env string) *config.Config {
	env = strings.ToUpper(env)
	return config.NewTestConfig(getDomain(env), getInternalAppSettings(env))
}

func getDomain(env string) constants.Domain {
	if env != "BOE" && env != "PRE" && env != "ONLINE" {
		panic("env must in [boe, pre, online]")
	}
	if env == "ONLINE" {
		return constants.DomainFeiShu
	}
	return constants.Domain(domainFeiShu(env))
}
