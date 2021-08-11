package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
)

func main() {
	// for redis store and logrus
	// var conf = sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
	// var conf = sample.TestConfig("https://open.feishu.cn")
	var conf = lark.NewInternalAppConfigByEnv(lark.DomainFeiShu)

	lark.WebHook.SetCardActionHandler(conf, func(ctx *lark.Context, cardAction *lark.CardAction) (interface{}, error) {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(cardAction.Action))
		return nil, nil
	})

	header := make(map[string][]string)
	// from http request header
	// and "X-Lark-Request-Timestamp"/"X-Lark-Request-Nonce"/"X-Lark-Signature" validate request, Required!

	// header["X-Request-Id"] = []string{"63278309j-yuewuyeu-7828389"}
	// header["X-Lark-Request-Timestamp"] = []string{"Monday, 09-Nov-20 23:33:53 CST"}
	// header["X-Lark-Request-Nonce"] = []string{"0404f57f-102e-4c91-bb32-a501ad0b7600"}
	// header["X-Lark-Signature"] = []string{"26cb59f4f5a91c4147d0xxxxxxxxxc4a36fb2c"}
	// header["X-Refresh-Token"] = []string{"acc4d5f2-4bc6-4394-a9d4-45e168fcde97"}

	req := &lark.HTTPRequest{
		Header: header,
		Body:   "", // from http request body
	}
	resp := lark.WebHook.CardRequestHandle(conf, req)
	fmt.Println(lark.Prettify(resp))
}
