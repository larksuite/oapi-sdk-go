package main

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/sample"
	"net/http"
	"os"
)

func main() {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("BOE_ISV_APP_ID"), os.Getenv("BOE_ISV_APP_SECRET"),
		os.Getenv("BOE_ISV_VERIFICATION_TOKEN"), os.Getenv("BOE_ISV_ENCRYPT_KEY")
	// BOE env
	larkApp := lark.NewApp("https://open.feishu.cn", appID, appSecret,
		lark.WithAppEventVerify(verificationToken, encryptKey),
		lark.WithAppType(lark.AppTypeMarketplace),            // marketplace app(ISV app)
		lark.WithLogger(sample.Logrus{}, lark.LogLevelDebug), // use logrus print log
		lark.WithStore(sample.NewRedisStore()),               // use redis store
	)

	marketplaceAppSendMessage(larkApp)

	// http server handle func
	// obtain app ticket
	http.HandleFunc("/webhook/event", func(writer http.ResponseWriter, request *http.Request) {
		rawRequest, err := lark.NewRawRequest(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		larkApp.Webhook.EventCommandHandle(context.Background(), rawRequest).Write(writer)
	})
	// startup http server
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}

func marketplaceAppSendMessage(larkApp *lark.App) {
	resp, err := larkApp.SendRequest(context.TODO(), "POST", "/open-apis/message/v4/send", map[string]interface{}{
		"user_id":  "beb9d1cc",
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": "test",
		},
	}, lark.AccessTokenTypeTenant, lark.WithTenantKey("13586be5aacf1748"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
