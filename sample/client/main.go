package main

import (
	"fmt"
	"net/http"
	"time"

	lark "github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
)

func createDefaultClient() {
	var feishu_client = lark.NewClient("appID", "appSecret")
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
		lark.WithAppType(larkcore.AppTypeSelfBuilt),
		lark.WithReqTimeout(3*time.Second),
		lark.WithDisableTokenCache(),
		lark.WithHelpdeskCredential("id", "token"),
		lark.WithLogger(larkcore.NewEventLogger()),
		lark.WithHttpClient(http.DefaultClient))
	fmt.Println(feishu_client)

}

func main() {

}
