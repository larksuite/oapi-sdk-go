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

package dispatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
)

type EventDispatcher struct {
	// 事件map,key为事件类型，value为事件处理器
	eventType2EventHandler map[string]larkevent.EventHandler
	// 事件回调签名token，消息解密key
	verificationToken string
	eventEncryptKey   string
	*larkcore.Config
}

func (dispatcher *EventDispatcher) Logger() larkcore.Logger {
	return dispatcher.Config.Logger
}

func (d *EventDispatcher) InitConfig(options ...larkevent.OptionFunc) {
	for _, option := range options {
		option(d.Config)
	}
	larkcore.NewLogger(d.Config)
}

func NewEventDispatcher(verificationToken, eventEncryptKey string) *EventDispatcher {
	reqDispatcher := &EventDispatcher{
		eventType2EventHandler: make(map[string]larkevent.EventHandler),
		verificationToken:      verificationToken,
		eventEncryptKey:        eventEncryptKey,
		Config:                 &larkcore.Config{Logger: larkcore.NewEventLogger()},
	}
	// 注册app_ticket事件
	reqDispatcher.eventType2EventHandler["app_ticket"] = &appTicketEventHandler{}
	return reqDispatcher
}

func processError(ctx context.Context, logger larkcore.Logger, path string, err error) *larkevent.EventResp {
	header := map[string][]string{}
	statusCode := http.StatusInternalServerError
	header[larkevent.ContentTypeHeader] = []string{larkevent.DefaultContentType}
	eventResp := &larkevent.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(larkevent.WebhookResponseFormat, err.Error())),
		StatusCode: statusCode,
	}
	logger.Error(ctx, fmt.Sprintf("handle event,path:%s,err: %v", path, err))
	return eventResp
}

func recoveryResult() *larkevent.EventResp {
	header := map[string][]string{}
	statusCode := http.StatusInternalServerError
	header[larkevent.ContentTypeHeader] = []string{larkevent.DefaultContentType}
	eventResp := &larkevent.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(larkevent.WebhookResponseFormat, "Server Internal Error")),
		StatusCode: statusCode,
	}
	return eventResp
}

func (d *EventDispatcher) Handle(ctx context.Context, req *larkevent.EventReq) (eventResp *larkevent.EventResp) {
	defer func() {
		if err := recover(); err != nil {
			d.Config.Logger.Error(ctx, fmt.Sprintf("handle event,path:%s, error:%v", req.RequestURI, err))
			eventResp = recoveryResult()
		}
	}()

	cipherEventJsonStr, err := d.ParseReq(ctx, req)
	if err != nil {
		return processError(ctx, d.Config.Logger, req.RequestURI, err)
	}

	plainEventJsonStr, err := d.DecryptEvent(ctx, cipherEventJsonStr)
	if err != nil {
		return processError(ctx, d.Config.Logger, req.RequestURI, err)
	}

	reqType, challenge, token, eventType, err := parse(plainEventJsonStr)
	if err != nil {
		return processError(ctx, d.Config.Logger, req.RequestURI, err)
	}
	if reqType != larkevent.ReqTypeChallenge {
		err = d.VerifySign(ctx, req)
		if err != nil {
			return processError(ctx, d.Config.Logger, req.RequestURI, err)
		}
	}

	result, err := d.DoHandle(ctx, reqType, eventType, challenge, token, plainEventJsonStr, req.RequestURI, req)
	if err != nil {
		return processError(ctx, d.Config.Logger, req.RequestURI, err)
	}
	return result
}

func (d *EventDispatcher) ParseReq(ctx context.Context, req *larkevent.EventReq) (string, error) {
	d.Config.Logger.Debug(ctx, fmt.Sprintf("event request: header:%v,body:%s", req.Header, string(req.Body)))
	if d.eventEncryptKey != "" {
		var encrypt larkevent.EventEncryptMsg
		err := json.Unmarshal(req.Body, &encrypt)
		if err != nil {
			err = fmt.Errorf("event message unmarshal failed:%v", err)
			return "", err
		}
		if encrypt.Encrypt == "" {
			err = fmt.Errorf("event  unmarshal failed,%s", "encrypted message is blank")
			return "", err
		}
		return encrypt.Encrypt, nil
	}
	return string(req.Body), nil
}

func (d *EventDispatcher) DecryptEvent(ctx context.Context, cipherEventJsonStr string) (str string, er error) {
	if d.eventEncryptKey != "" {
		body, err := larkevent.EventDecrypt(cipherEventJsonStr, d.eventEncryptKey)
		if err != nil {
			err = fmt.Errorf("event message decryption failed:%v", err)
			return "", err
		}
		return string(body), nil
	}
	return cipherEventJsonStr, nil
}

func (d *EventDispatcher) VerifySign(ctx context.Context, req *larkevent.EventReq) error {
	if d.eventEncryptKey == "" {
		return nil
	}
	// 解析签名头
	requestTimestamps := req.Header[larkevent.EventRequestTimestamp]
	requestNonces := req.Header[larkevent.EventRequestNonce]
	var requestTimestamp = ""
	var requestNonce = ""
	if len(requestTimestamps) > 0 {
		requestTimestamp = requestTimestamps[0]
	}
	if len(requestNonces) > 0 {
		requestNonce = requestNonces[0]
	}
	// 执行sha256签名计算
	targetSign := larkevent.Signature(requestTimestamp, requestNonce,
		d.eventEncryptKey, string(req.Body))

	sourceSigns := req.Header[larkevent.EventSignature]
	var sourceSign = ""
	if len(sourceSigns) > 0 {
		sourceSign = sourceSigns[0]
	}

	// 验签
	if targetSign == sourceSign {
		return nil
	}
	return errors.New("the result of signature verification failed")
}

func parse(plainEventJsonStr string) (larkevent.ReqType, string, string, string, error) {
	fuzzy := &larkevent.EventFuzzy{}
	err := json.Unmarshal([]byte(plainEventJsonStr), fuzzy)
	if err != nil {
		err = fmt.Errorf("event json unmarshal, err: %v", err)
		return "", "", "", "", err
	}
	if fuzzy.Encrypt != "" {
		err = errors.New("event data is encrypted, Need to set up the `EncryptKey` for your APP")
		return "", "", "", "", err
	}

	reqType := larkevent.ReqType(fuzzy.Type)
	var eventType string
	token := fuzzy.Token
	challenge := fuzzy.Challenge
	if fuzzy.Event != nil {
		if et, ok := fuzzy.Event.Type.(string); ok {
			eventType = et
		}
	}
	if fuzzy.Header != nil {
		token = fuzzy.Header.Token
		eventType = fuzzy.Header.EventType
	}

	return reqType, challenge, token, eventType, nil
}

func (d *EventDispatcher) getErrorResp(ctx context.Context, header map[string][]string, err error) *larkevent.EventResp {
	eventResp := &larkevent.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(larkevent.WebhookResponseFormat, err.Error())),
		StatusCode: http.StatusInternalServerError,
	}
	d.Config.Logger.Error(ctx, fmt.Sprintf("event handle err: %v", err))
	return eventResp
}

func (d *EventDispatcher) AuthByChallenge(ctx context.Context, reqType larkevent.ReqType, challenge, token string) (*larkevent.EventResp, error) {
	if reqType == larkevent.ReqTypeChallenge {
		if token != d.verificationToken {
			err := errors.New("the result of auth by challenge failed")
			return nil, err
		}

		header := map[string][]string{}
		header[larkevent.ContentTypeHeader] = []string{larkevent.DefaultContentType}
		eventResp := larkevent.EventResp{
			Header:     header,
			Body:       []byte(fmt.Sprintf(larkevent.ChallengeResponseFormat, challenge)),
			StatusCode: http.StatusOK,
		}
		d.Config.Logger.Info(ctx, fmt.Sprintf("AuthByChallenge Success"))
		return &eventResp, nil
	}
	return nil, nil
}

func (d *EventDispatcher) DoHandle(ctx context.Context, reqType larkevent.ReqType, eventType, challenge, token,
	plainEventJsonStr string, path string, req *larkevent.EventReq) (*larkevent.EventResp, error) {
	// auth by challenge
	resp, err := d.AuthByChallenge(ctx, reqType, challenge, token)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return resp, nil
	}

	// 查找处理器
	handler := d.eventType2EventHandler[eventType]
	if handler == nil {
		err = &notFoundEventHandlerErr{eventType: eventType}
		header := map[string][]string{}
		header[larkevent.ContentTypeHeader] = []string{larkevent.DefaultContentType}
		eventResp := &larkevent.EventResp{
			Header:     header,
			Body:       []byte(fmt.Sprintf(larkevent.WebhookResponseFormat, err.Error())),
			StatusCode: http.StatusOK,
		}
		d.Config.Logger.Error(ctx, fmt.Sprintf("handle event,path:%s,error:%v", path, err.Error()))
		return eventResp, nil
	}

	// 反序列化
	eventMsg := handler.Event()
	if _, ok := handler.(*defaultHandler); !ok {
		err = json.Unmarshal([]byte(plainEventJsonStr), eventMsg)
		if err != nil {
			return nil, err
		}
	} else {
		eventMsg = req
	}

	if msg, ok := eventMsg.(larkevent.EventHandlerModel); ok {
		msg.RawReq(req)
	}

	// 执行处理器
	err = handler.Handle(ctx, eventMsg)
	if err != nil {
		return nil, err
	}

	//返回结果
	header := map[string][]string{}
	header[larkevent.ContentTypeHeader] = []string{larkevent.DefaultContentType}
	eventResp := &larkevent.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(larkevent.WebhookResponseFormat, "success")),
		StatusCode: http.StatusOK,
	}
	return eventResp, nil
}

type notFoundEventHandlerErr struct {
	eventType string
}

func (e notFoundEventHandlerErr) Error() string {
	return fmt.Sprintf("event type: %s, not found handler", e.eventType)
}

type defaultHandler struct {
	handler func(context.Context, *larkevent.EventReq) error
}

func (h *defaultHandler) Event() interface{} {
	return nil
}

func (h *defaultHandler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*larkevent.EventReq))
}
