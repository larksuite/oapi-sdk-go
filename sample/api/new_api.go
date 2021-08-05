package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	im "github.com/larksuite/oapi-sdk-go/service/im/v1"
)

var newConf = lark.NewInternalAppConfigByEnv(lark.DomainFeiShu)

func main() {
	newConf.SetLogLevel(lark.LogLevelDebug)
	testSendCardMessage1()
	testMessageCreate1()
}

// send card message
func testSendCardMessage1() {
	coreCtx := lark.WrapContext(context.Background())
	body := map[string]interface{}{
		"user_id":  "77bbc392",
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": "test",
		},
	}
	ret := make(map[string]interface{})
	req := lark.NewAPIRequest("/open-apis/message/v4/send", "POST",
		lark.AccessTokenTypeTenant, body, &ret,
		//应用市场应用 lark.SetTenantKey("TenantKey"),
	)
	err := lark.SendAPIRequest(coreCtx, newConf, req)
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(ret))
}

func testMessageCreate1() {
	coreCtx := lark.WrapContext(context.Background())
	var imService = im.NewService(newConf)
	reqCall := imService.Messages.Create(coreCtx, &im.MessageCreateReqBody{
		// ReceiveId: "b1g6b445",
		ReceiveId: "ou_a11d2bcc7d852afbcaf37e5b3ad01f7e",
		Content:   "{\"text\":\"<at user_id=\\\"ou_a11d2bcc7d852afbcaf37e5b3ad01f7e\\\">Tom</at> test content\"}",
		MsgType:   "text",
	})
	reqCall.SetReceiveIdType("open_id")
	message, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	fmt.Println(coreCtx.GetHeader())
	if err != nil {
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(message))
}
