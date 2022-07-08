package main

import (
	"context"
	"fmt"
	"os"

	client "github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	larkim "github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func rawApi1() {
	// 历史base url 需要在httppath上拼接上baseurl
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = client.NewClient(appID, appSecret,
		client.WithLogLevel(core.LogLevelDebug),
		client.WithLogReqRespInfoAtDebugLevel(true))

	resp, err := cli.Post(context.Background(), "https://www.feishu.cn/approval/openapi/v2/approval/get", map[string]interface{}{
		"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
	}, core.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}

func rawApi2() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = client.NewClient(appID, appSecret,
		client.WithLogLevel(core.LogLevelDebug),
		client.WithLogReqRespInfoAtDebugLevel(true))

	content := larkim.NewTextMsgBuilder().
		Text("加多").
		Line().
		TextLine("hello").
		TextLine("world").
		AtUser("ou_c245b0a7dff2725cfa2fb104f8b48b9d", "陆续").
		Build()

	resp, err := cli.Post(context.Background(), "/open-apis/im/v1/messages?receive_id_type=open_id", map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    "{\"text\":\"<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at> test content\"}",
	}, core.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err, content)
		return
	}

	fmt.Println(resp)
}

func main() {
	rawApi1()
}
