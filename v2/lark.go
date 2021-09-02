package lark

import (
	"encoding/base64"
	"fmt"
)

func NewApp(domain Domain, appID, appSecret string, options ...AppOptionFunc) *App {
	app := &App{
		domain:   domain,
		settings: &appSettings{type_: AppTypeCustom},
		logger:   newLoggerProxy(LogLevelError, nil),
		store:    &defaultStore{},
	}
	options = append(options, withAppCredential(appID, appSecret))
	for _, optionFunc := range options {
		optionFunc(app)
	}
	app.Webhook = newWebhook(app)
	return app
}

func withAppCredential(appID, appSecret string) AppOptionFunc {
	return func(app *App) {
		app.settings.id = appID
		app.settings.secret = appSecret
	}
}

func WithAppEventVerify(verificationToken, encryptKey string) AppOptionFunc {
	return func(app *App) {
		app.settings.verificationToken = verificationToken
		app.settings.eventEncryptKey = encryptKey
	}
}

func WithAppHelpdeskCredential(helpdeskID, helpdeskToken string) AppOptionFunc {
	return func(app *App) {
		app.settings.helpDeskID = helpdeskID
		app.settings.helpDeskToken = helpdeskToken
		if helpdeskID != "" && helpdeskToken != "" {
			app.settings.helpdeskAuthToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", helpdeskID, helpdeskToken)))
		}
	}
}

func WithLogger(logger logger, level LogLevel) AppOptionFunc {
	return func(lark *App) {
		lark.logger = newLoggerProxy(level, logger)
	}
}

func WithAppType(appType AppType) AppOptionFunc {
	return func(app *App) {
		app.settings.type_ = appType
	}
}

func WithStore(store store) AppOptionFunc {
	return func(app *App) {
		app.store = store
	}
}

type AppOptionFunc func(*App)

type App struct {
	domain   Domain
	settings *appSettings
	logger   logger
	store    store
	Webhook  *webhook
}

type appSettings struct {
	type_  AppType
	id     string
	secret string

	verificationToken string
	eventEncryptKey   string

	helpDeskID        string
	helpDeskToken     string
	helpdeskAuthToken string
}

type webhook struct {
	app                    *App
	actionHandler          cardActionHandler
	eventType2EventHandler map[string]eventHandler
}

func newWebhook(app *App) *webhook {
	return &webhook{
		app:                    app,
		eventType2EventHandler: map[string]eventHandler{},
	}
}
