package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	appID, appSecret, verificationToken, encryptKey := os.Getenv("APP_ID"), os.Getenv("APP_SECRET"),
		os.Getenv("VERIFICATION_TOKEN"), os.Getenv("ENCRYPT_KEY")

	larkApp := lark.NewApp(lark.DomainFeiShu,
		lark.WithAppCredential(appID, appSecret),
		lark.WithEventVerify(verificationToken, encryptKey),
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	lark.Webhook.CardActionHandler(larkApp, func(ctx context.Context, request *lark.RawRequest,
		action *lark.CardAction) (interface{}, error) {
		fmt.Println(request)
		fmt.Println(lark.Prettify(action))
		return nil, nil
	})

	// http server handle func
	http.HandleFunc("/webhook/card", func(writer http.ResponseWriter, request *http.Request) {
		rawBody, _ := ioutil.ReadAll(request.Body)
		resp := lark.Webhook.CardActionHandle(context.Background(), larkApp, &lark.RawRequest{
			Header:  request.Header,
			RawBody: rawBody,
		})
		writer.WriteHeader(resp.StatusCode)
		for k, vs := range resp.Header {
			for _, v := range vs {
				writer.Header().Set(k, v)
			}
		}
		writer.Write(resp.RawBody)
	})
	// startup http server
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
