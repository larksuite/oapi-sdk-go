package event

import (
	"context"
	"fmt"
	"net/http"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func main() {
	//1212121212
	handler := dispatcher.NewEventDispatcher("", "1212121212")
	handler.OnCustomizedEvent("custom_event_type", func(ctx context.Context, event *larkevent.EventReq) error {
		// 原生消息体
		fmt.Println(string(event.Body))
		fmt.Println(larkcore.Prettify(event.Header))
		fmt.Println(larkcore.Prettify(event.RequestURI))
		fmt.Println(event.RequestId())

		// 处理消息
		cipherEventJsonStr, err := handler.ParseReq(ctx, event)
		if err != nil {
			//  错误处理
			return err
		}

		plainEventJsonStr, err := handler.DecryptEvent(ctx, cipherEventJsonStr)
		if err != nil {
			//  错误处理
			return err
		}

		// 处理解密后的 消息体
		fmt.Println(plainEventJsonStr)

		return nil
	}).OnP2CardNewProtocalURLPreviewGet(func(ctx context.Context, event *dispatcher.URLPreviewGetEvent) (*dispatcher.URLPreviewGetResponse, error) {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return &dispatcher.URLPreviewGetResponse{
			Inline: &dispatcher.Inline{
				Title:    "title",
				ImageKey: "image_key",
			},
			Card: &dispatcher.Card{
				Type: "raw",
				Data: &larkcard.MessageCard{},
			},
		}, nil
	}).OnP2CardNewProtocalCardActionTrigger(func(ctx context.Context, event *dispatcher.CardActionTriggerEvent) (*dispatcher.CardActionTriggerReponse, error) {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return &dispatcher.CardActionTriggerReponse{
			Toast: &dispatcher.Toast{
				Type:    "info",
				Content: "toast",
			},
			Card: &dispatcher.Card{
				Type: "raw",
				Data: &larkcard.MessageCard{},
			},
		}, nil
	})

	// 注册 http 路由
	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler,
		larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动服务
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}
