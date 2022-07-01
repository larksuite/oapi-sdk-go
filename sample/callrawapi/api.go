package main

import (
	"context"
	"fmt"

	client "github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
)

func main() {
	config := &core.Config{
		Domain:    client.FeishuDomain,
		AppId:     "appId",
		AppSecret: "appSecret",
	}

	resp, err := core.SendGet(context.Background(), config, "/open-apis/message/v4/send", core.AccessTokenTypeUser, map[string]interface{}{
		"user_id":  "77bbc392",
		"msg_type": "text",
		"content":  "content",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
