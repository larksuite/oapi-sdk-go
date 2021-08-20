package lark

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var once sync.Once

func (h webhook) EventCommandHandle(ctx context.Context, app *App, req *RawRequest) *RawResponse {
	return eventCommandHandle(ctx, app, req)
}

func eventCommandHandle(ctx context.Context, app *App, request *RawRequest) *RawResponse {
	once.Do(func() {
		if app.settings.type_ == AppTypeMarketplace {
			EventHandler(app, "app_ticket", &appTicketEventHandler{app: app})
		}
	})
	httpEvent := &httpEvent{
		request:  request,
		response: &RawResponse{},
	}
	httpEvent.do(ctx, app)
	return httpEvent.response
}

type eventHandler interface {
	Event() interface{}
	Handle(context.Context, *RawRequest, interface{}) error
}

var appId2Type2EventHandler = map[string]map[string]eventHandler{}

func EventHandler(app *App, eventType string, handler eventHandler) {
	appID := app.settings.id
	type2EventHandler, ok := appId2Type2EventHandler[appID]
	if !ok {
		type2EventHandler = map[string]eventHandler{}
		appId2Type2EventHandler[appID] = type2EventHandler
	}
	type2EventHandler[eventType] = handler
}

func EventHandleFunc(app *App, eventType string, handler func(context.Context, *RawRequest) error) {
	EventHandler(app, eventType, &defaultHandler{handler: handler})
}

type defaultHandler struct {
	handler func(context.Context, *RawRequest) error
}

func (h *defaultHandler) Event() interface{} {
	return nil
}

func (h *defaultHandler) Handle(ctx context.Context, req *RawRequest, event interface{}) error {
	return h.handler(ctx, req)
}

type httpEvent struct {
	request  *RawRequest
	response *RawResponse
}

func (e *httpEvent) do(ctx context.Context, app *App) {
	var type_ webhookType
	var token string
	var eventType string
	var challenge string
	var err error
	defer func() {
		e.response.StatusCode = http.StatusOK
		e.response.Header = map[string][]string{}
		e.response.Header.Set(contentTypeHeader, defaultContentType)
		if err != nil {
			if _, ok := err.(*notFoundEventHandlerErr); ok {
				e.response.RawBody = []byte(fmt.Sprintf(webhookResponseFormat, err.Error()))
				return
			}
			app.logger.Error(ctx, fmt.Sprintf("event handle err: %v", err))
			e.response.StatusCode = http.StatusInternalServerError
			e.response.RawBody = []byte(fmt.Sprintf(webhookResponseFormat, err.Error()))
			return
		}
		if type_ == webhookTypeChallenge {
			e.response.RawBody = []byte(fmt.Sprintf(challengeResponseFormat, challenge))
			return
		}
		e.response.RawBody = []byte(fmt.Sprintf(webhookResponseFormat, "success"))
		return
	}()
	var body = e.request.RawBody
	app.logger.Debug(ctx, fmt.Sprintf("event request: %v", e.request))
	if app.settings.encryptKey != "" {
		var encrypt eventAESMsg
		err = json.Unmarshal(e.request.RawBody, &encrypt)
		if err != nil {
			err = fmt.Errorf("event json unmarshal, err:%v", err)
			return
		}
		body, err = eventDecrypt(encrypt.Encrypt, app.settings.encryptKey)
		if err != nil {
			err = fmt.Errorf("event decrypt, err:%v", err)
			return
		}
		app.logger.Debug(ctx, fmt.Sprintf("event decrypt: %s", string(body)))
	}
	fuzzy := &eventFuzzy{}
	err = json.Unmarshal(body, fuzzy)
	if err != nil {
		err = fmt.Errorf("event json unmarshal, err: %v", err)
		return
	}
	type_ = webhookType(fuzzy.Type)
	token = fuzzy.Token
	challenge = fuzzy.Challenge
	if fuzzy.Event != nil {
		if et, ok := fuzzy.Event.Type.(string); ok {
			eventType = et
		}
	}
	if fuzzy.Header != nil {
		token = fuzzy.Header.Token
		eventType = fuzzy.Header.EventType
	}
	if type_ == webhookTypeChallenge {
		if token != app.settings.verificationToken {
			err = errors.New("event token not equal app settings token")
			return
		}
		return
	}
	var handler eventHandler
	if eventType2EventHandler, ok := appId2Type2EventHandler[app.settings.id]; ok {
		if h, ok := eventType2EventHandler[eventType]; ok {
			handler = h
		}
	}
	if handler == nil {
		err = &notFoundEventHandlerErr{eventType: eventType}
		return
	}
	event := handler.Event()
	err = json.Unmarshal(body, event)
	if err != nil {
		return
	}
	err = handler.Handle(ctx, e.request, event)
}

// eventDecrypt returns decrypt bytes
func eventDecrypt(encrypt string, secret string) ([]byte, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return nil, newDecryptErr(fmt.Sprintf("base64 decode error: %v", err))
	}
	if len(buf) < aes.BlockSize {
		return nil, newDecryptErr("cipher too short")
	}
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:sha256.Size])
	if err != nil {
		return nil, newDecryptErr(fmt.Sprintf("AES new cipher error %v", err))
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return nil, newDecryptErr("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return buf[n : m+1], nil
}

type EventBase struct {
	Ts    string `json:"ts"`
	UUID  string `json:"uuid"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

type eventFuzzy struct {
	Schema    string       `json:"schema"`
	Token     string       `json:"token"`
	Type      string       `json:"type"`
	Challenge string       `json:"challenge"`
	Header    *EventHeader `json:"header"`
	Event     *struct {
		Type interface{} `json:"type"`
	} `json:"event"`
}

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

type eventAESMsg struct {
	Encrypt string `json:"encrypt"`
}

type notFoundEventHandlerErr struct {
	eventType string
}

func (e notFoundEventHandlerErr) Error() string {
	return fmt.Sprintf("event type: %s, not found handler", e.eventType)
}

type appTicketEventData struct {
	AppID     string `json:"app_id"`
	Type      string `json:"type"`
	TenantKey string `json:"tenant_key"`
	AppTicket string `json:"app_ticket"`
}

type appTicketEvent struct {
	*EventBase
	Event *appTicketEventData `json:"event"`
}

type appTicketEventHandler struct {
	app   *App
	event *appTicketEvent
}

func (h *appTicketEventHandler) Event() interface{} {
	h.event = &appTicketEvent{}
	return h.event
}

func (h *appTicketEventHandler) Handle(ctx context.Context, req *RawRequest, event interface{}) error {
	appTicketEvent := event.(*appTicketEvent)
	return h.app.store.Put(ctx, fmt.Sprintf("%s-%s", appTicketKeyPrefix, appTicketEvent.Event.AppID),
		appTicketEvent.Event.AppTicket, time.Hour*12)
}
