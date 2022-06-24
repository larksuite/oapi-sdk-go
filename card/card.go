package card

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
)

type CardActionHandler struct {
	verificationToken string
	eventEncryptKey   string
	handler           func(context.Context, *CardAction) (interface{}, error)
	*core.Config
}

func processError(ctx context.Context, logger core.Logger, err error) *event.EventResp {
	header := map[string][]string{}
	statusCode := http.StatusInternalServerError
	header[event.ContentTypeHeader] = []string{event.DefaultContentType}
	eventResp := &event.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(event.WebhookResponseFormat, err.Error())),
		StatusCode: statusCode,
	}
	logger.Error(ctx, fmt.Sprintf("card handle err: %v", err))
	return eventResp
}

func (h *CardActionHandler) Handle(ctx context.Context, req *event.EventReq) *event.EventResp {
	h.Config.Logger.Debug(ctx, fmt.Sprintf("card request: header:%v,body:%s", req.Header, string(req.Body)))

	cardAction := &CardAction{}
	err := json.Unmarshal(req.Body, cardAction)
	if err != nil {
		return processError(ctx, h.Config.Logger, err)
	}

	if event.ReqType(cardAction.Type) != event.ReqTypeChallenge {
		err = h.VerifySign(ctx, req)
		if err != nil {
			return processError(ctx, h.Config.Logger, err)
		}
	}

	result, err := h.DoHandle(ctx, cardAction)
	if err != nil {
		return processError(ctx, h.Config.Logger, err)
	}
	return result
}

func (h *CardActionHandler) Logger() core.Logger {
	return h.Config.Logger
}

func (h *CardActionHandler) InitConfig(options ...event.OptionFunc) {
	for _, option := range options {
		option(h.Config)
	}
	core.NewLogger(h.Config)
}

func NewCardActionHandler(verificationToken, eventEncryptKey string, handler func(context.Context, *CardAction) (interface{}, error)) *CardActionHandler {
	h := &CardActionHandler{
		verificationToken: verificationToken,
		eventEncryptKey:   eventEncryptKey,
		handler:           handler,
		Config:            &core.Config{Logger: core.NewEventLogger()},
	}
	return h
}

func (h *CardActionHandler) Event() interface{} {
	return &CardAction{}
}

var notFoundCardHandlerErr = errors.New("card action handler not found")

func (h *CardActionHandler) AuthByChallenge(ctx context.Context, cardAction *CardAction) (*event.EventResp, error) {
	header := map[string][]string{}
	header[event.ContentTypeHeader] = []string{event.DefaultContentType}
	hookType := event.ReqType(cardAction.Type)
	challenge := cardAction.Challenge
	if hookType == event.ReqTypeChallenge {
		if h.verificationToken != cardAction.Token {
			err := errors.New("the result of auth by challenge failed")
			return nil, err
		}
		eventResp := event.EventResp{
			Header:     header,
			Body:       []byte(fmt.Sprintf(event.ChallengeResponseFormat, challenge)),
			StatusCode: http.StatusOK,
		}
		return &eventResp, nil
	}
	return nil, nil
}
func (h *CardActionHandler) DoHandle(ctx context.Context, cardAction *CardAction) (*event.EventResp, error) {
	var err error
	// auth by challenge
	resp, err := h.AuthByChallenge(ctx, cardAction)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return resp, nil
	}

	// 校验行为执行器
	handler := h.handler
	if handler == nil {
		err = notFoundCardHandlerErr
		return nil, err
	}

	// 执行事件行为处理器
	result, err := handler(ctx, cardAction)
	if err != nil {
		return nil, err
	}

	header := map[string][]string{}
	header[event.ContentTypeHeader] = []string{event.DefaultContentType}
	if result == nil {
		eventResp := &event.EventResp{
			Header:     header,
			Body:       []byte(fmt.Sprintf(event.WebhookResponseFormat, "success")),
			StatusCode: http.StatusOK,
		}
		return eventResp, nil
	}

	var respBody []byte
	switch r := result.(type) {
	case string:
		respBody = []byte(r)
	case *CustomResp:
		statusCode := r.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		b, err := json.Marshal(r.Body)
		if err != nil {
			return nil, err
		}

		eventResp := &event.EventResp{
			Header:     header,
			Body:       b,
			StatusCode: statusCode,
		}
		return eventResp, nil
	default:
		respBody, err = json.Marshal(result)
	}

	eventResp := &event.EventResp{
		Header:     header,
		Body:       respBody,
		StatusCode: http.StatusOK,
	}

	return eventResp, err
}

func (h *CardActionHandler) VerifySign(ctx context.Context, req *event.EventReq) error {
	if h.verificationToken == "" {
		return nil
	}

	// 解析签名头
	requestTimestamps := req.Header[event.EventRequestTimestamp]
	requestNonces := req.Header[event.EventRequestNonce]

	var requestTimestamp = ""
	var requestNonce = ""
	if len(requestTimestamps) > 0 {
		requestTimestamp = requestTimestamps[0]
	}
	if len(requestNonces) > 0 {
		requestNonce = requestNonces[0]
	}

	// 执行sha1签名计算
	targetSign := Signature(requestTimestamp, requestNonce,
		h.verificationToken,
		string(req.Body))

	sourceSigns := req.Header[event.EventSignature]
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

func Signature(timestamp, nonce, token, body string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(token)
	b.WriteString(body)
	bs := []byte(b.String())
	h := sha1.New()
	_, _ = h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
