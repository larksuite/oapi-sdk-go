package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
)

func mockCardAction() []byte {
	var token = "v"
	value := map[string]interface{}{}
	value["value"] = "sdfsfd"
	value["tag"] = "button"

	cardAction := &larkcard.CardAction{
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

	cardActionBody := &larkcard.CardActionBody{
		CardAction: cardAction,
		Challenge:  "121212",
		Type:       "url_verification1",
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
	var token = "v"

	//var b = "{\"open_id\":\"ou_d840b2e2be16b3e0091bc0c79220e1fa\",\"user_id\":\"16fd348g\",\"open_message_id\":\"om_dce5707d696ee4952ebedaf1ee762ed2\",\"tenant_key\":\"736588c9260f175d\",\"token\":\"v\",\"action\":{\"value\":{\"key\":\"value\"},\"tag\":\"button\"}}"
	sourceSign := larkcard.Signature(timestamp, nonce, token, string(body))
	//fmt.Println(sourceSign)
	// 添加header
	req.Header.Set(larkevent.EventRequestTimestamp, timestamp)
	req.Header.Set(larkevent.EventRequestNonce, nonce)
	req.Header.Set(larkevent.EventSignature, sourceSign)
	req.Header.Set("X-Tt-Logid", "logid111111111111111")

	// 模拟推送卡片消息
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 结果处理
	fmt.Println(resp.StatusCode)
	fmt.Println(larkcore.Prettify(resp.Header))
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
