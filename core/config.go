package larkcore

import (
	"time"
)

type Config struct {
	BaseUrl                    string
	AppId                      string
	AppSecret                  string
	HelpDeskId                 string
	HelpDeskToken              string
	HelpdeskAuthToken          string
	ReqTimeout                 time.Duration
	LogLevel                   LogLevel
	HttpClient                 HttpClient
	Logger                     Logger
	AppType                    AppType
	EnableTokenCache           bool
	TokenCache                 Cache
	LogReqRespInfoAtDebugLevel bool
}
