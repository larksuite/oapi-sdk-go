package config

import (
	"github.com/larksuite/oapi-sdk-go/core/constants"
)

type Settings struct {
}

type AppSettings struct {
	AppType   constants.AppType
	AppID     string
	AppSecret string

	VerificationToken string
	EncryptKey        string
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
