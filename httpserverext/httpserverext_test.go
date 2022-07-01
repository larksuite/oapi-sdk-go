package httpserverext

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func TestStartHttpServer(t *testing.T) {
	// 创建card处理器
	cardHandler := card.NewCardActionHandler("12", "", func(ctx context.Context, cardAction *card.CardAction) (interface{}, error) {
		fmt.Println(core.Prettify(cardAction))
		return nil, nil
	})

	// 创建事件处理器
	handler := dispatcher.NewEventDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *larkim.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnMessageReadV1(func(ctx context.Context, event *larkim.MessageMessageReadEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	// 注册事件 和 卡片路径
	http.HandleFunc("/webhook/event", NewEventHandlerFunc(handler))
	http.HandleFunc("/webhook/card", NewCardActionHandlerFunc(cardHandler))

	// 启动服务
	//err := http.ListenAndServe(":9999", nil)
	//if err != nil {
	//	panic(err)
	//}
}

func mockRequest() *http.Request {
	var token = "12"
	value := map[string]interface{}{}
	value["value"] = "sdfsfd"
	value["tag"] = "button"

	cardAction := &card.CardAction{
		OpenID:        "ou_sdfimx9948345",
		UserID:        "eu_sd923r0sdf5",
		OpenMessageID: "om_abcdefg1234567890",
		TenantKey:     "d32004232",
		Token:         token,
		Action: &struct {
			Value    map[string]interface{} `json:"value"`
			Tag      string                 `json:"tag"`
			Option   string                 `json:"option"`
			Timezone string                 `json:"timezone"`
		}{
			Value: value,
			Tag:   "button",
		},
	}

	cardActionBody := &card.CardActionBody{
		CardAction: cardAction,
		Challenge:  "121212",
		Type:       "url_verification",
	}

	body, _ := json.Marshal(cardActionBody)
	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(body))
	req.Header.Set("key1", "value1")
	req.Header.Set("key2", "value2")
	return req
}

func TestTranslate(t *testing.T) {
	req := mockRequest()
	eventReq, err := translate(context.Background(), req)
	if err != nil {
		t.Errorf("translate failed ,%v", err)
		return
	}

	fmt.Println(core.Prettify(eventReq.Header))
	fmt.Println(string(eventReq.Body))
}
