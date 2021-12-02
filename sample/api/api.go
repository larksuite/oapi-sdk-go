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
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	ctx := context.Background()
	//sendMessage(ctx, larkApp)
	//uploadImage(ctx, larkApp)
	//downloadImage(ctx, larkApp)
	getJsTicket(ctx, larkApp)
}

func sendMessage(ctx context.Context, larkApp *lark.App) {
	resp, err := larkApp.SendRequest(ctx, http.MethodGet, "/open-apis/message/v4/send",
		lark.AccessTokenTypeTenant, map[string]interface{}{
			"user_id":  "77bbc392",
			"msg_type": "text",
			"content":  &lark.MessageText{Text: "test"},
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", resp.RequestId())
	fmt.Println(resp)
	fmt.Println()
	fmt.Println()
}

func getJsTicket(ctx context.Context, larkApp *lark.App) {
	resp, err := larkApp.SendRequest(ctx, http.MethodGet, "/open-apis/jssdk/ticket/get",
		lark.AccessTokenTypeApp, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", resp.RequestId())
	fmt.Println(resp)
	fmt.Println()
	fmt.Println()
}

func uploadImage(ctx context.Context, larkApp *lark.App) {
	img, err := os.Open("test.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer img.Close()
	resp, err := larkApp.SendRequest(ctx, http.MethodPost, "/open-apis/image/v4/put",
		lark.AccessTokenTypeTenant, lark.NewFormdata().AddField("image_type", "message").AddFile("image", img))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", resp.RequestId())
	fmt.Println(resp)
	fmt.Println()
	fmt.Println()
}

func downloadImage(ctx context.Context, larkApp *lark.App) {
	resp, err := larkApp.SendRequest(ctx, http.MethodGet, "/open-apis/image/v4/get?image_key=img_v2_b49b2df3-f277-467b-a650-d352839f4b6g",
		lark.AccessTokenTypeTenant, nil, lark.WithFileDownload())
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
		return
	}
	err = ioutil.WriteFile("test_download_v2.png", resp.RawBody, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", resp.RequestId())
	fmt.Println(resp)
}
