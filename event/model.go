/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkevent

import (
	"net/http"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type EventHeader struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
}

type EventV1Header struct {
	AppID     string `json:"app_id"`
	OpenAppID string `json:"open_chat_id"`
	OpenID    string `json:"open_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`
}

type EventV2Base struct {
	Schema string       `json:"schema"`
	Header *EventHeader `json:"header"`
}

func (base *EventV2Base) TenantKey() string {
	if base != nil && base.Header != nil {
		return base.Header.TenantKey
	}
	return ""
}

type EventV2Body struct {
	EventV2Base
	Challenge string      `json:"challenge"`
	Event     interface{} `json:"event"`
	Type      string      `json:"type"`
}

type EventReq struct {
	Header     map[string][]string
	Body       []byte
	RequestURI string
}

func (req *EventReq) RequestId() string {
	logID := req.Header[larkcore.HttpHeaderKeyLogId]
	if len(logID) > 0 {
		return logID[0]
	}
	logID = req.Header[larkcore.HttpHeaderKeyRequestId]
	if len(logID) > 0 {
		return logID[0]
	}
	return ""
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
