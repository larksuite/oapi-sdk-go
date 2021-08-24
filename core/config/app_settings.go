package config

import (
	"encoding/base64"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"os"
)

type AppSettingsOpt func(*AppSettings)

type AppSettings struct {
	AppType constants.AppType

	AppID     string
	AppSecret string

	VerificationToken string
	EncryptKey        string

	HelpDeskID    string
	HelpDeskToken string
}

func (s *AppSettings) helpDeskAuthorization() string {
	if s.HelpDeskID != "" && s.HelpDeskToken != "" {
		helpdeskAuthToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.HelpDeskID, s.HelpDeskToken)))
		return helpdeskAuthToken
	}
	return ""
}

func GetISVAppSettingsByEnv() *AppSettings {
	return NewISVAppSettingsByOpts(getAppSettingsOptsByEnv()...)
}

func GetInternalAppSettingsByEnv() *AppSettings {
	return NewInternalAppSettingsByOpts(getAppSettingsOptsByEnv()...)
}

func NewISVAppSettings(appID, appSecret, verificationToken, encryptKey string) *AppSettings {
	return NewISVAppSettingsByOpts(SetAppCredentials(appID, appSecret), SetAppEventKey(verificationToken, encryptKey))
}

func NewInternalAppSettings(appID, appSecret, verificationToken, encryptKey string) *AppSettings {
	return NewInternalAppSettingsByOpts(SetAppCredentials(appID, appSecret), SetAppEventKey(verificationToken, encryptKey))
}

func NewISVAppSettingsByOpts(opts ...AppSettingsOpt) *AppSettings {
	return newAppSettingsByOpts(constants.AppTypeISV, opts...)
}

func NewInternalAppSettingsByOpts(opts ...AppSettingsOpt) *AppSettings {
	return newAppSettingsByOpts(constants.AppTypeInternal, opts...)
}

func newAppSettingsByOpts(appType constants.AppType, optFns ...AppSettingsOpt) *AppSettings {
	settings := &AppSettings{AppType: appType}
	for _, opt := range optFns {
		opt(settings)
	}
	if settings.AppID == "" || settings.AppSecret == "" {
		panic("appID or appSecret is empty")
	}
	return settings
}

func getAppSettingsOptsByEnv() []AppSettingsOpt {
	var opts []AppSettingsOpt
	appID, appSecret, verificationToken, encryptKey, helpDeskID, helpDeskToken := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY"), os.Getenv("HELP_DESK_ID"), os.Getenv("HELP_DESK_TOKEN")
	opts = append(opts, SetAppCredentials(appID, appSecret))
	opts = append(opts, SetAppEventKey(verificationToken, encryptKey))
	opts = append(opts, SetHelpDeskCredentials(helpDeskID, helpDeskToken))
	return opts
}

func SetAppCredentials(appID, appSecret string) AppSettingsOpt {
	return func(settings *AppSettings) {
		settings.AppID = appID
		settings.AppSecret = appSecret
	}
}

func SetAppEventKey(verificationToken, encryptKey string) AppSettingsOpt {
	return func(settings *AppSettings) {
		settings.VerificationToken = verificationToken
		settings.EncryptKey = encryptKey
	}
}

func SetHelpDeskCredentials(helpDeskID, helpDeskToken string) AppSettingsOpt {
	return func(settings *AppSettings) {
		settings.HelpDeskID = helpDeskID
		settings.HelpDeskToken = helpDeskToken
	}
}
