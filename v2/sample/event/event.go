package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"net/http"
	"os"
)

func main() {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY")

	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithAppEventVerify(verificationToken, encryptKey),
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	// @robot message handle
	larkApp.Webhook.EventHandleFunc("im.message.receive_v1", func(ctx context.Context, req *lark.RawRequest) error {
		fmt.Println(req.RequestId())
		fmt.Println(req)
		return nil
	})

	// http server handle func
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
