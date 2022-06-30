package core

import (
	"net/http"
	"time"
)

type Config struct {
	Domain                     string
	AppId                      string
	AppSecret                  string
	HelpDeskId                 string
	HelpDeskToken              string
	HelpdeskAuthToken          string
	ReqTimeout                 time.Duration
	LogLevel                   LogLevel
	HttpClient                 *http.Client
	Logger                     Logger
	AppType                    AppType
	EnableTokenCache           bool
	TokenCache                 Cache
	LogReqRespInfoAtDebugLevel bool
}
