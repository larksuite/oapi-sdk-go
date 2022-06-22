package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/httpserverext"
)

func getCard() *card.CardAction {
	// 返回卡片消息
	value := map[string]interface{}{}
	value["value"] = "sdfsfd"
	value["tag"] = "button"

	result := &card.CardAction{
		OpenID:        "ou_sdfimx9948345",
		UserID:        "eu_sd923r0sdf5",
		OpenMessageID: "om_abcdefg1234567890",
		TenantKey:     "d32004232",
		Token:         "121",
		Action: &struct {
			Value    map[string]interface{} `json:"value"`
			Tag      string                 `json:"tag"`
			Option   string                 `json:"option"`
			Timezone string                 `json:"timezone"`
		}{
			Value: value,
			Tag:   "button",
		},
	}
	fmt.Println(result)
	return result
}

func getCustomResp() interface{} {
	toastBody := card.CustomToastBody{
		Content: "sfsfsdfsd",
		I18n: &card.I18n{
			ZhCn: "ZhCn",
			EnCn: "EnCn",
			JaJp: "JaJp",
		},
	}
	body, _ := json.Marshal(toastBody)

	resp := card.CustomResp{
		StatusCode: 400,
		Body:       body,
	}
	return &resp
}
func main() {

	// 创建card处理器
	cardHandler := card.NewCardActionHandler("12", "", func(ctx context.Context, cardAction *card.CardAction) (interface{}, error) {
		fmt.Println(core.Prettify(cardAction))

		// 返回卡片消息
		//return getCard(), nil

		//custom resp
		return getCustomResp(), nil

		// 无返回值
		return nil, nil
	})

	// 注册处理器
	http.HandleFunc("/webhook/card", httpserverext.NewCardActionHandlerFunc(cardHandler, event.WithLogLevel(core.LogLevelDebug)))

	// 开发者启动服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
