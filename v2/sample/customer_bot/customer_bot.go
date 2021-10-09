package main

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v2"
)

func main() {
	customerBot := lark.NewCustomerBot("https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxxxx", "u8TfAYQYZWvWRubw6EKQUe")
	resp, err := customerBot.SendMessage(context.TODO(), "text", &lark.MessageText{Text: "test"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("request id: %s \n", resp.RequestId())
	fmt.Println(lark.Prettify(resp))
}
