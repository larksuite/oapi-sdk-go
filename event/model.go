package event

import (
	"net/http"
)

type EventHeader struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
}

type EventV2Base struct {
	Schema string       `json:"schema"`
	Header *EventHeader `json:"header"`
}

type EventV2Body struct {
	EventV2Base
	Challenge string      `json:"challenge"`
	Event     interface{} `json:"event"`
	Type      string      `json:"type"`
}

type EventReq struct {
	Header map[string][]string
	Body   []byte
}

type EventResp struct {
	Header     http.Header
	Body       []byte
	StatusCode int
}

type EventBase struct {
	Ts    string `json:"ts"`
	UUID  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

type EventEncryptMsg struct {
	Encrypt string `json:"encrypt"`
}

type EventFuzzy struct {
	Encrypt   string       `json:"encrypt"`
	Schema    string       `json:"schema"`
	Token     string       `json:"token"`
	Type      string       `json:"type"`
	Challenge string       `json:"challenge"`
	Header    *EventHeader `json:"header"`
	Event     *struct {
		Type interface{} `json:"type"`
	} `json:"event"`
}

const (
	EventRequestNonce     = "X-Lark-Request-Nonce"
	EventRequestTimestamp = "X-Lark-Request-Timestamp"
	EventSignature        = "X-Lark-Signature"
)

type ReqType string

const (
	ReqTypeChallenge     ReqType = "url_verification"
	ReqTypeEventCallBack ReqType = "event_callback"
)

const userAgentHeader = "User-Agent"
const ContentTypeHeader = "Content-Type"
const ContentTypeJson = "application/json"
const DefaultContentType = ContentTypeJson + "; charset=utf-8"
const WebhookResponseFormat = `{"msg":"%s"}`
const ChallengeResponseFormat = `{"challenge":"%s"}`
