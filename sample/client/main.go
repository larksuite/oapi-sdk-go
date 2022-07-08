package main

import (
	"fmt"
	"net/http"
	"time"

	client "github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
)

func createDefaultClient() {
	var feishu_client = client.NewClient("appID", "appSecret")
	fmt.Println(feishu_client)
}

func createClientWithLogLevel() {
	var feishu_client = client.NewClient("appID", "appSecret", client.WithLogLevel(core.LogLevelDebug))
	fmt.Println(feishu_client)
}

func createClientWithAllOptions() {
	var feishu_client = client.NewClient("appID", "appSecret",
		client.WithLogLevel(core.LogLevelDebug),
		client.WithOpenBaseUrl(client.FeishuBaseUrl),
		client.WithAppType(core.AppTypeSelfBuilt),
		client.WithReqTimeout(3*time.Second),
		client.WithDisableTokenCache(),
		client.WithHelpdeskCredential("id", "token"),
		client.WithLogger(core.NewEventLogger()),
		client.WithHttpClient(http.DefaultClient))
	fmt.Println(feishu_client)

}

func main() {

}
