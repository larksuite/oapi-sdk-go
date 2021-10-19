package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/sample"
	lark "github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/service/im/v1"
	"net/http"
	"os"
)

func main() {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("BOE_ISV_APP_ID"), os.Getenv("BOE_ISV_APP_SECRET"),
		os.Getenv("BOE_ISV_VERIFICATION_TOKEN"), os.Getenv("BOE_ISV_ENCRYPT_KEY")
	// BOE env
	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithAppEventVerify(verificationToken, encryptKey),
		lark.WithAppType(lark.AppTypeMarketplace),            // marketplace app(ISV app)
		lark.WithLogger(sample.Logrus{}, lark.LogLevelDebug), // use logrus print log
		lark.WithStore(sample.NewRedisStore()),               // use redis store
	)

	marketplaceAppSendMessage(context.Background(), larkApp)

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

func marketplaceAppSendMessage(ctx context.Context, larkApp *lark.App) {
	messageText := &lark.MessageText{Text: "Tom test content"}
	content, err := messageText.JSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	messageCreateResp, err := im.New(larkApp).Messages.Create(ctx, &im.MessageCreateReq{
		ReceiveIdType: lark.StringPtr("user_id"),
		Body: &im.MessageCreateReqBody{
			ReceiveId: lark.StringPtr("beb9d1cc"),
			MsgType:   lark.StringPtr("text"),
			Content:   lark.StringPtr(content),
		},
	}, lark.WithTenantKey("13586be5aacf1748")) // setting TenantKey
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", messageCreateResp.RequestId())
	if messageCreateResp.Code != 0 {
		fmt.Println(messageCreateResp.CodeError)
		return
	}
	fmt.Println(lark.Prettify(messageCreateResp.Data))
	fmt.Println()
	fmt.Println()
}
