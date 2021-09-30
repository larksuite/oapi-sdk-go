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

	// card action handler
	// return new card
	larkApp.Webhook.CardActionHandler(func(ctx context.Context, request *lark.RawRequest,
		action *lark.CardAction) (interface{}, error) {
		fmt.Println(request)
		fmt.Println(lark.Prettify(action))

		card := &lark.MessageCard{
			Config: &lark.MessageCardConfig{WideScreenMode: lark.BoolPtr(true)},
			Elements: []lark.MessageCardElement{&lark.MessageCardMarkdown{
				Content: "**test**",
			}},
		}
		return card, nil
	})

	// http server handle func
	http.HandleFunc("/webhook/card", func(writer http.ResponseWriter, request *http.Request) {
		rawRequest, err := lark.NewRawRequest(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		larkApp.Webhook.CardActionHandle(context.Background(), rawRequest).Write(writer)
	})
	// startup http server
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
