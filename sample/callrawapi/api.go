package main

import (
	"context"
	"fmt"

	client "github.com/feishu/oapi-sdk-go"
	"github.com/feishu/oapi-sdk-go/core"
)

func main() {
	config := &core.Config{
		Domain:    client.FeishuDomain,
		AppId:     "appId",
		AppSecret: "appSecret",
	}

	resp, err := core.SendRequest(context.Background(), config, "get", "/open-apis/message/v4/send", []core.AccessTokenType{core.AccessTokenTypeUser}, map[string]interface{}{
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
