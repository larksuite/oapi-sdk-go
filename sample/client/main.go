package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func createDefaultClient() {
	var feishu_client = lark.NewClient("appID", "appSecret",
		lark.WithMarketplaceApp())
	fmt.Println(feishu_client)
}

func createClientWithLogLevel() {
	var feishu_client = lark.NewClient("appID", "appSecret",
		lark.WithLogLevel(larkcore.LogLevelDebug))
	fmt.Println(feishu_client)
}

func createClientWithAllOptions() {
	var feishu_client = lark.NewClient("appID", "appSecret",
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithOpenBaseUrl(lark.FeishuBaseUrl),
		lark.WithMarketplaceApp(),
		lark.WithReqTimeout(3*time.Second),
		lark.WithEnableTokenCache(false),
		lark.WithHelpdeskCredential("id", "token"),
		lark.WithLogger(larkcore.NewEventLogger()),
		lark.WithHttpClient(http.DefaultClient),
		lark.WithLogReqRespInfoAtDebugLevel(true))
	fmt.Println(feishu_client)

}

func main() {

}
