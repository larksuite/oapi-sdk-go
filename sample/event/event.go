package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/httpserverext"
	"github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func main() {

	//1212121212
	handler := dispatcher.NewEventDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *larkim.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnMessageReadV1(func(ctx context.Context, event *larkim.MessageReadEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnUserCreatedV3(func(ctx context.Context, event *larkcontact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, event.WithLogLevel(core.LogLevelDebug)))

	// 开发者启动服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
