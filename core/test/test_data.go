package test

import (
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"os"
)

func domainFeiShuStaging() string {
	return os.Getenv("DomainFeiShuStaging")
}

func getISVAppSettings() *config.AppSettings {
	appID, appSecret, verificationToken, eventEncryptKey := os.Getenv("ISVAppID"),
		os.Getenv("ISVAppSecret"), os.Getenv("ISVVerificationToken"), os.Getenv("ISVEventEncryptKey")
	return config.NewISVAppSettings(appID, appSecret, verificationToken, eventEncryptKey)
}

func GetStagingISVConf() *config.Config {
	return config.NewTestConfig(constants.Domain(domainFeiShuStaging()), getISVAppSettings())
}

func GetStagingInternalConf() *config.Config {
	return config.NewTestConfig(constants.Domain(domainFeiShuStaging()), getInternalAppSettings())
}

func getInternalAppSettings() *config.AppSettings {
	appID, appSecret, verificationToken, eventEncryptKey := os.Getenv("InternalAppID"),
		os.Getenv("InternalAppSecret"), os.Getenv("InternalVerificationToken"), os.Getenv("InternalEventEncryptKey")
	return config.NewInternalAppSettings(appID, appSecret, verificationToken, eventEncryptKey)
}

func GetOnlineInternalConf() *config.Config {
	return config.NewTestConfig(constants.DomainFeiShu, getInternalAppSettings())
}
