package card

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/feishu/oapi-sdk-go/core"
	"github.com/feishu/oapi-sdk-go/event"
)

type CardActionHandler struct {
	verificationToken string
	eventEncryptKey   string
	handler           func(context.Context, *CardAction) (interface{}, error)
	event.ReqHandler
	*core.Config
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

var notFoundCardHandlerErr = errors.New("card not found handler")

func (h *CardActionHandler) VerifyUrl(ctx context.Context, plainEventJsonStr string) (*event.EventResp, error) {
	// 反序列化
	challengeMsg := &cardChallenge{}
	err := json.Unmarshal([]byte(plainEventJsonStr), challengeMsg)
	if err != nil {
		return nil, err
	}

	//URL验证
	header := map[string][]string{}
	header[event.ContentTypeHeader] = []string{event.DefaultContentType}
	hookType := event.WebhookType(challengeMsg.Type)
	challenge := challengeMsg.Challenge
	if hookType == event.WebhookTypeChallenge {
		if h.verificationToken != challengeMsg.Token {
			err = errors.New("card challenge token not equal settings token")
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
func (h *CardActionHandler) DoHandle(ctx context.Context, plainEventJsonStr string) (*event.EventResp, error) {
	// 校验行为执行器
	var err error
	handler := h.handler
	if handler == nil {
		err = notFoundCardHandlerErr
		return nil, err
	}
	// 反序列化事件内容
	cardAction := &CardAction{}
	err = json.Unmarshal([]byte(plainEventJsonStr), cardAction)
	if err != nil {
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
		eventResp := &event.EventResp{
			Header:     header,
			Body:       r.Body,
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

func (h *CardActionHandler) ParseReq(ctx context.Context, req *event.EventReq) (string, error) {
	h.Config.Logger.Debug(ctx, fmt.Sprintf("cardAction request: %v", req))
	return string(req.Body), nil
}

func (h *CardActionHandler) DecryptEvent(ctx context.Context, cipherEventJsonStr string) (string, error) {
	return cipherEventJsonStr, nil
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
	return errors.New("signature error")
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
