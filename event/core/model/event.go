package model

import (
	"github.com/larksuite/oapi-sdk-go/core"
)

const Version1 = "1.0"
const Version2 = "2.0"

type HTTPEvent struct {
	Request   *core.OapiRequest
	Input     []byte
	Response  *core.OapiResponse
	Schema    string
	Type      string
	EventType string
	Challenge string
	Err       error
}

type BaseEvent struct {
	Ts    string `json:"ts"`
	UUID  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

type BaseEventData struct {
	AppID     string `json:"app_id"`
	Type      string `json:"type"`
	TenantKey string `json:"tenant_key"`
}

type Fuzzy struct {
	Schema    string  `json:"schema"`
	Token     string  `json:"token"`
	Type      string  `json:"type"`
	Challenge string  `json:"challenge"`
	Header    *Header `json:"header"`
	Event     *struct {
		Type string `json:"type"`
	} `json:"event"`
}

type Header struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
}

type BaseEventV2 struct {
	Schema string  `json:"schema"`
	Header *Header `json:"header"`
}
