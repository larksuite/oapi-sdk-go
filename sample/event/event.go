package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/feishu/oapi-sdk-go/card"
	"github.com/feishu/oapi-sdk-go/core"
	"github.com/feishu/oapi-sdk-go/dispatcher"
	"github.com/feishu/oapi-sdk-go/httpserverext"
	"github.com/feishu/oapi-sdk-go/service/contact/v3"
	"github.com/feishu/oapi-sdk-go/service/im/v1"
)

type CardActionBody struct {
	*card.CardAction
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func main() {

	handler := dispatcher.NewEventReqDispatcher("v", "1212121212").MessageReceiveV1(func(ctx context.Context, event *im.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).MessageMessageReadV1(func(ctx context.Context, event *im.MessageMessageReadEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).UserCreatedV3(func(ctx context.Context, event *contact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	http.HandleFunc("/webhook/event", httpserverext.NewEventReqHandlerFunc(handler))

	// 开发者启动服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
