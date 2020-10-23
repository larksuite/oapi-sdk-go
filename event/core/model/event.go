package model

import (
	"net/http"
)

const Version1 = "v1"
const Version2 = "v2"

type HTTPEvent struct {
	HTTPRequest  *http.Request
	Input        []byte
	HTTPResponse http.ResponseWriter
	Version      string
	Type         string
	EventType    string
	Challenge    string
	Err          error
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

type NotData struct {
	Version   string  `json:"version"`
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
	Version string  `json:"version"`
	Header  *Header `json:"header"`
}
