package lark

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"net/http"
)

var WebHook webHook = webHook{}

type webHook struct{}

type HTTPRequest struct {
	Header http.Header
	Body   string
}

type HTTPResponse = core.OapiResponse
