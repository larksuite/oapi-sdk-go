package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/channel"
	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

func main() {
	// 支持命令行参数
	var appID string
	var appSecret string
	flag.StringVar(&appID, "app_id", "", "Lark App ID (也可以通过 APP_ID 环境变量传入)")
	flag.StringVar(&appSecret, "app_secret", "", "Lark App Secret (也可以通过 APP_SECRET 环境变量传入)")
	flag.Parse()

	// 从环境变量获取配置作 fallback
	if appID == "" {
		appID = os.Getenv("APP_ID")
	}
	if appSecret == "" {
		appSecret = os.Getenv("APP_SECRET")
	}

	if appID == "" || appSecret == "" {
		fmt.Println("启动失败: 缺少必要的参数")
		fmt.Println("使用方式: go run main.go -app_id=cli_xxx -app_secret=xxx")
		fmt.Println("或者设置环境变量: APP_ID=cli_xxx APP_SECRET=xxx go run main.go")
		return
	}

	// 1. 创建 Lark API Client
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))

	// 2. 创建 EventDispatcher，用于 WebSocket 的事件路由
	eventHandler := dispatcher.NewEventDispatcher("", "")

	// 3. 创建 WebSocket Client，并将 EventHandler 注入
	wsClient := larkws.NewClient(appID, appSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)

	// 4. 创建 Channel 抽象实例，将 client 和 wsClient 传入
	ch := channel.NewChannel(client, wsClient)

	// 5. 注册消息处理逻辑
	ch.OnMessage(func(ctx context.Context, msg *types.NormalizedMessage) error {
		fmt.Printf("收到来自 %s 的消息: %s\n", msg.UserID, msg.Content)

		if msg.Content == "stream" {
			// 测试 Stream (流式响应)
			stream, err := ch.Stream(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Markdown:  "正在思考中...",
			})
			if err != nil {
				return err
			}

			// 模拟慢速生成文本
			for i := 1; i <= 5; i++ {
				stream.Append(ctx, fmt.Sprintf("这是第 %d 行流式文本...\n", i))
				time.Sleep(1 * time.Second)
			}

			stream.Append(ctx, "\n**完成！**")
			stream.Close(ctx)

		} else {
			// 测试普通 Send (发送消息)
			_, err := ch.Send(ctx, &types.SendInput{
				ReceiveID: msg.ChatID,
				Text:      "Echo: " + msg.Content,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	// 6. 注册卡片交互处理逻辑 (可选)
	ch.OnCardAction(func(ctx context.Context, msg *types.NormalizedMessage) error {
		fmt.Printf("收到来自 %s 的卡片交互，Value: %v\n", msg.UserID, msg.Content)
		return nil
	})

	// 7. 启动 WebSocket Client，开始监听事件
	fmt.Println("正在启动 WebSocket Client...")
	err := wsClient.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
