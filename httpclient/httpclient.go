package httpclient

import (
	"net/http"

	"github.com/feishu/oapi-sdk-go/core"
)

func NewHttpClient(config *core.Config) *http.Client {
	if config.HttpClient == nil {
		config.HttpClient = &http.Client{Timeout: config.ReqTimeout}
	}

	return config.HttpClient

}
