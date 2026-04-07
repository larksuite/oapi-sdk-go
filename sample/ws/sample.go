package main

import (
	"context"
	"fmt"
	"os"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

type envWSClientAssertionProvider struct{}

func (p *envWSClientAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	return &larkcore.Token{Value: os.Getenv("CLIENT_ASSERTION")}, nil
}

func main() {
	// 注册事件回调
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			fmt.Printf("[ OnP2MessageReceiveV1 access ], data: %s\n", larkcore.Prettify(event))
			return nil
		}).
		OnCustomizedEvent("这里填入你要自定义订阅的 event 的 key,例如 out_approval", func(ctx context.Context, event *larkevent.EventReq) error {
			fmt.Printf("[ OnCustomizedEvent access ], type: message, data: %s\n", string(event.Body))
			return nil
		})

	// 创建Client
	cli := larkws.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)

	// 建立长连接
	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}
}

func startWithClientAssertion(ctx context.Context) error {
	cli := larkws.NewClient(os.Getenv("APP_ID"), "",
		larkws.WithClientAssertionProvider(&envWSClientAssertionProvider{}),
	)
	return cli.Start(ctx)
}
