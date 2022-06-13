package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/feishu/oapi-sdk-go/card"
	"github.com/feishu/oapi-sdk-go/core"
	"github.com/feishu/oapi-sdk-go/event"
)

func mockCardAction() []byte {
	var token = "12"
	value := map[string]interface{}{}
	value["value"] = "sdfsfd"
	value["tag"] = "button"
	cardAction := &card.CardAction{
		OpenID:        "ou_sdfimx9948345",
		UserID:        "eu_sd923r0sdf5",
		OpenMessageID: "om_abcdefg1234567890",
		TenantKey:     "d32004232",
		Token:         token,
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

	cardActionBody := &card.CardActionBody{
		CardAction: cardAction,
		Challenge:  "121212",
		Type:       "url_verification",
	}

	body, _ := json.Marshal(cardActionBody)
	return body
}
func main() {

	//mock body
	body := mockCardAction()

	// 创建http req
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9999/webhook/card", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 计算签名
	var timestamp = "timestamp"
	var nonce = "nonce"
	var token = "token"
	sourceSign := card.Signature(timestamp, nonce, token, string(body))

	// 添加header
	req.Header.Set(event.EventRequestTimestamp, timestamp)
	req.Header.Set(event.EventRequestNonce, nonce)
	req.Header.Set(event.EventSignature, sourceSign)

	// 模拟推送卡片消息
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 结果处理
	fmt.Println(resp.StatusCode)
	fmt.Println(core.Prettify(resp.Header))
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
