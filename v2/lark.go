package lark

import (
	"encoding/base64"
	"fmt"
)

func NewApp(domain Domain, options ...AppOptionFunc) *App {
	app := &App{
		domain:   domain,
		settings: &AppSettings{type_: AppTypeCustom},
		logger:   NewLoggerProxy(LogLevelError, nil),
		store:    &defaultStore{},
	}
	for _, optionFunc := range options {
		optionFunc(app)
	}
	return app
}

func WithAppCredential(appID, appSecret string) AppOptionFunc {
	return func(app *App) {
		app.settings.id = appID
		app.settings.secret = appSecret
	}
}

func WithEventVerify(verificationToken, encryptKey string) AppOptionFunc {
	return func(app *App) {
		app.settings.verificationToken = verificationToken
		app.settings.encryptKey = encryptKey
	}
}

func WithHelpdeskCredential(helpdeskID, helpdeskToken string) AppOptionFunc {
	return func(app *App) {
		app.settings.helpDeskID = helpdeskID
		app.settings.helpDeskToken = helpdeskToken
		if helpdeskID != "" && helpdeskToken != "" {
			app.settings.helpdeskAuthToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", helpdeskID, helpdeskToken)))
		}
	}
}

func WithLogger(logger Logger, level LogLevel) AppOptionFunc {
	return func(lark *App) {
		lark.logger = NewLoggerProxy(level, logger)
	}
}

func WithAppType(appType AppType) AppOptionFunc {
	return func(app *App) {
		app.settings.type_ = appType
	}
}

func WithStore(store Store) AppOptionFunc {
	return func(app *App) {
		app.store = store
	}
}

type AppOptionFunc func(*App)

type App struct {
	domain   Domain
	settings *AppSettings
	logger   Logger
	store    Store
}

type AppSettings struct {
	type_  AppType
	id     string
	secret string

	verificationToken string
	encryptKey        string

	helpDeskID        string
	helpDeskToken     string
	helpdeskAuthToken string
}

var Webhook = webhook{}

type webhook struct {
}
