package dispatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/feishu/oapi-sdk-go/core"
	"github.com/feishu/oapi-sdk-go/event"
)

type EventReqDispatcher struct {
	// 事件map,key为事件类型，value为事件处理器
	eventType2EventHandler map[string]event.EventHandler
	// 事件回调签名token，消息解密key
	verificationToken string
	eventEncryptKey   string
	event.ReqHandler
	*core.Config
}

func NewEventReqDispatcher(verificationToken, eventEncryptKey string) *EventReqDispatcher {
	reqDispatcher := &EventReqDispatcher{
		eventType2EventHandler: make(map[string]event.EventHandler),
		verificationToken:      verificationToken,
		eventEncryptKey:        eventEncryptKey,
		Config:                 &core.Config{Logger: core.NewEventLogger()},
	}
	// 注册app_ticket事件
	reqDispatcher.eventType2EventHandler["app_ticket"] = &appTicketEventHandler{}

	return reqDispatcher
}

func (d *EventReqDispatcher) ParseReq(ctx context.Context, req *event.EventReq) (string, error) {
	d.Config.Logger.Debug(ctx, fmt.Sprintf("event request: %v", req))
	if d.eventEncryptKey != "" {
		var encrypt event.EventEncryptMsg
		err := json.Unmarshal(req.Body, &encrypt)
		if err != nil {
			err = fmt.Errorf("event json unmarshal, err:%v", err)
			return "", err
		}

		if encrypt.Encrypt == "" {
			err = fmt.Errorf("event json unmarshal error,%v", "encrypt content is blank")
			return "", err
		}

		return encrypt.Encrypt, nil

	}
	return string(req.Body), nil
}

func (d *EventReqDispatcher) DecryptEvent(ctx context.Context, cipherEventJsonStr string) (string, error) {
	if d.eventEncryptKey != "" {
		body, err := event.EventDecrypt(cipherEventJsonStr, d.eventEncryptKey)
		if err != nil {
			err = fmt.Errorf("event decrypt, err:%v", err)
			return "", err
		}
		return string(body), nil
	}
	return cipherEventJsonStr, nil
}

func (d *EventReqDispatcher) VerifySign(ctx context.Context, req *event.EventReq) error {
	if d.eventEncryptKey == "" {
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

	// 执行sha256签名计算
	targetSign := event.Signature(requestTimestamp, requestNonce,
		d.eventEncryptKey, string(req.Body))

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

func parse(plainEventJsonStr string) (event.WebhookType, string, string, string, error) {
	fuzzy := &event.EventFuzzy{}
	err := json.Unmarshal([]byte(plainEventJsonStr), fuzzy)
	if err != nil {
		err = fmt.Errorf("event json unmarshal, err: %v", err)
		return "", "", "", "", err
	}
	if fuzzy.Encrypt != "" {
		err = errors.New("event data is encrypted, Need to set up the `EncryptKey` for your APP")
		return "", "", "", "", err
	}

	hookType := event.WebhookType(fuzzy.Type)
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

	return hookType, challenge, token, eventType, nil
}

func (d *EventReqDispatcher) getErrorResp(ctx context.Context, header map[string][]string, err error) *event.EventResp {
	eventResp := &event.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(event.WebhookResponseFormat, err.Error())),
		StatusCode: http.StatusInternalServerError,
	}
	d.Config.Logger.Error(ctx, fmt.Sprintf("event handle err: %v", err))
	return eventResp
}

func (d *EventReqDispatcher) VerifyUrl(ctx context.Context, plainEventJsonStr string) (*event.EventResp, error) {
	// 解析数据
	header := map[string][]string{}
	header[event.ContentTypeHeader] = []string{event.DefaultContentType}
	hookType, challenge, token, _, err := parse(plainEventJsonStr)
	if err != nil {
		return nil, err
	}

	// 处理url验证
	if hookType == event.WebhookTypeChallenge {
		if token != d.verificationToken {
			err = errors.New("event token not equal settings token")
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

func (d *EventReqDispatcher) DoHandle(ctx context.Context, plainEventJsonStr string) (*event.EventResp, error) {
	// 解析数据
	header := map[string][]string{}
	header[event.ContentTypeHeader] = []string{event.DefaultContentType}
	_, _, _, eventType, err := parse(plainEventJsonStr)
	if err != nil {
		return nil, err
	}

	// 查找处理器
	handler := d.eventType2EventHandler[eventType]
	if handler == nil {
		err = &notFoundEventHandlerErr{eventType: eventType}
		return nil, err
	}

	// 反序列化
	eventMsg := handler.Event()
	if _, ok := handler.(*defaultHandler); !ok {
		err = json.Unmarshal([]byte(plainEventJsonStr), eventMsg)
		if err != nil {
			return nil, err
		}
	}

	// 执行处理器
	err = handler.Handle(ctx, eventMsg)
	if err != nil {
		return nil, err
	}

	//返回结果
	eventResp := &event.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(event.WebhookResponseFormat, "success")),
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
	handler func(context.Context) error
}

func (h *defaultHandler) Event() interface{} {
	return nil
}

func (h *defaultHandler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx)
}
