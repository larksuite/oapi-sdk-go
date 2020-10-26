package config

import (
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"os"
)

type AppSettings struct {
	AppType   constants.AppType
	AppID     string
	AppSecret string

	VerificationToken string
	EncryptKey        string
}

func GetISVAppSettingsByEnv() *AppSettings {
	appID, appSecret, verificationToken, encryptKey := getAppSettingsByEnv()
	return NewISVAppSettings(appID, appSecret, verificationToken, encryptKey)
}

func GetInternalAppSettingsByEnv() *AppSettings {
	appID, appSecret, verificationToken, encryptKey := getAppSettingsByEnv()
	return NewInternalAppSettings(appID, appSecret, verificationToken, encryptKey)
}

func NewISVAppSettings(appID, appSecret, verificationToken, eventEncryptKey string) *AppSettings {
	return newAppSettings(constants.AppTypeISV, appID, appSecret, verificationToken, eventEncryptKey)
}

func NewInternalAppSettings(appID, appSecret, verificationToken, eventEncryptKey string) *AppSettings {
	return newAppSettings(constants.AppTypeInternal, appID, appSecret, verificationToken, eventEncryptKey)
}

func newAppSettings(appType constants.AppType, appID, appSecret, verificationToken, eventEncryptKey string) *AppSettings {
	if appID == "" || appSecret == "" {
		panic("appID or appSecret is empty")
	}
	return &AppSettings{
		AppType:           appType,
		AppID:             appID,
		AppSecret:         appSecret,
		VerificationToken: verificationToken,
		EncryptKey:        eventEncryptKey,
	}
}

func getAppSettingsByEnv() (string, string, string, string) {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY")
	if appID == "" {
		panic("environment variables not exist `APP_ID`")
	}
	if appSecret == "" {
		panic("environment variables not exist `APP_SECRET`")
	}
	return appID, appSecret, verificationToken, encryptKey
}
