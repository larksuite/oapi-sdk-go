package main

import (
	"context"
	"fmt"
	"os"

	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func rawApi1() {
	// 历史base url 需要在httppath上拼接上baseurl
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithLogReqRespInfoAtDebugLevel(true))

	resp, err := cli.Post(context.Background(), "https://www.feishu.cn/approval/openapi/v2/approval/get", map[string]interface{}{
		"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
	}, larkcore.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode) // http status code
	fmt.Println(resp.Header)     // http header
	fmt.Println(resp.RawBody)    // http body
}

func rawApi2() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithLogReqRespInfoAtDebugLevel(true))

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
	}, larkcore.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err, content)
		return
	}

	fmt.Println(resp)
}

func main() {
	rawApi2()
}
