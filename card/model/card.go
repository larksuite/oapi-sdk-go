package model

import (
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"net/http"
)

const (
	LarkRequestTimestamp    = "X-Lark-Request-Timestamp"
	LarkRequestRequestNonce = "X-Lark-Request-Nonce"
	LarkSignature           = "X-Lark-Signature"
	LarkRefreshToken        = "X-Refresh-Token"
)

type Header struct {
	Timestamp    string
	Nonce        string
	Signature    string
	RefreshToken string
	RequestID    string
}

type HTTPCard struct {
	Header       *Header
	HTTPRequest  *http.Request
	Input        []byte
	HTTPResponse http.ResponseWriter
	Type         constants.CallbackType
	Output       interface{}
	Challenge    string
	Err          error
}

type Challenge struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type Base struct {
	OpenID        string `json:"open_id"`
	UserID        string `json:"user_id"`
	OpenMessageID string `json:"open_message_id"`
	TenantKey     string `json:"tenant_key"`
	Token         string `json:"token"`
	Timezone      string `json:"timezone"`
}

type Card struct {
	*Base
	Action *Action `json:"action"`
}

type Action struct {
	Value    map[string]string `json:"value"`
	Tag      string            `json:"tag"`
	Option   string            `json:"option"`
	Timezone string            `json:"timezone"`
}
