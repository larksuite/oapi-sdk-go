package core

import (
	"time"

	"github.com/larksuite/oapi-sdk-go/httpclient"
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
	HttpClient                 httpclient.HttpClient
	Logger                     Logger
	AppType                    AppType
	EnableTokenCache           bool
	TokenCache                 Cache
	LogReqRespInfoAtDebugLevel bool
}
