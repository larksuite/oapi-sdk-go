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
	sendMessage(ctx, larkApp)
	//uploadImage(ctx, larkApp)
	//downloadImage(ctx, larkApp)
}

func sendMessage(ctx context.Context, larkApp *lark.App) {
	resp, err := larkApp.SendRequest(ctx, http.MethodGet, "/open-apis/message/v4/send", lark.AccessTokenTypeTenant, map[string]interface{}{
		"user_id":  "77bbc392",
		"msg_type": "text",
		"content":  `{"config":{"enable_forward":true},"header":{"template":"blue","title":{"content":"Header title","tag":"plain_text"}},"elements":[{"alt":{"content":"img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg","tag":"plain_text"},"compact_width":false,"custom_width":300,"img_key":"img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg","tag":"img","title":{"content":"img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg","tag":"plain_text"}},{"actions":[{"confirm":{"title":{"content":"Title","tag":"plain_text"},"text":{"content":"Text","tag":"plain_text"}},"tag":"button","text":{"content":"button","tag":"plain_text"},"type":"danger","value":{"value":"1"}}],"layout":"flow","tag":"action"},{"content":"**Markdown**","tag":"markdown"},{"extra":{"confirm":{"title":{"content":"Title","tag":"plain_text"},"text":{"content":"Text","tag":"plain_text"}},"tag":"button","text":{"content":"button","tag":"plain_text"},"type":"danger","value":{"value":"1"}},"tag":"div","text":{"content":"text","tag":"plain_text"}}],"card_link":{"url":"https://www.feishu.cn","android_url":"https://www.feishu.cn","ios_url":"https://www.feishu.cn","pc_url":"https://www.feishu.cn"}}`,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.RequestId())
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
	fmt.Println(resp.RequestId())
	fmt.Println(resp)
	fmt.Println()
	fmt.Println()
}

func downloadImage(ctx context.Context, larkApp *lark.App) {
	resp, err := larkApp.SendRequest(ctx, http.MethodGet, "/open-apis/image/v4/get?image_key=img_v2_a0cea146-64d2-4dcb-94c7-636586fea98g",
		lark.AccessTokenTypeTenant, nil)
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
	fmt.Println(resp.RequestId())
	fmt.Println(resp)
}
