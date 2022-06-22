package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/service/contact/v3"
)

type EventV2Body struct {
	contact.UserCreatedEvent        // 什么类型消息就替换为什么
	Challenge                string `json:"challenge"`
	Type                     string `json:"type"`
}

func mockEncryptedBody(encrypteKey string) []byte {

	userEvent := &contact.UserEvent{
		OpenId:     core.StringPtr("ou_7dab8a3d3cdcc9da365777c7ad535d62"),
		UnionId:    core.StringPtr("on_576833b917gda3d939b9a3c2d53e72c8"),
		UserId:     core.StringPtr("e33ggbyz"),
		Name:       core.StringPtr("张三"),
		EmployeeNo: core.StringPtr("employee_no"),
	}

	usersCreatedEvent := contact.UserCreatedEvent{
		EventV2Base: &event.EventV2Base{
			Schema: "2.0",
			Header: &event.EventHeader{
				EventID:    "f7984f25108f8137722bb63cee927e66",
				EventType:  "contact.user.created_v3",
				AppID:      "cli_xxxxxxxx",
				TenantKey:  "xxxxxxx",
				CreateTime: "1603977298000000",
				Token:      "v",
			},
		},
		Event: &contact.UserCreatedEventData{Object: userEvent},
	}

	eventBody := EventV2Body{
		UserCreatedEvent: usersCreatedEvent,
		Challenge:        "1212",
		Type:             "url_verification1",
	}

	en, _ := core.EncryptedEventMsg(context.Background(), eventBody, encrypteKey)
	fmt.Println(encrypteKey)

	encrypt := event.EventEncryptMsg{Encrypt: en}
	body1, _ := json.Marshal(encrypt)

	return body1
}

func mockEvent() []byte {
	userEvent := &contact.UserEvent{
		OpenId:     core.StringPtr("ou_7dab8a3d3cdcc9da365777c7ad535d62"),
		UnionId:    core.StringPtr("on_576833b917gda3d939b9a3c2d53e72c8"),
		UserId:     core.StringPtr("e33ggbyz"),
		Name:       core.StringPtr("张三"),
		EmployeeNo: core.StringPtr("employee_no"),
	}

	usersCreatedEvent := contact.UserCreatedEvent{
		EventV2Base: &event.EventV2Base{
			Schema: "2.0",
			Header: &event.EventHeader{
				EventID:    "f7984f25108f8137722bb63cee927e66",
				EventType:  "contact.user.created_v3",
				AppID:      "cli_xxxxxxxx",
				TenantKey:  "xxxxxxx",
				CreateTime: "1603977298000000",
				Token:      "v",
			},
		},
		Event: &contact.UserCreatedEventData{Object: userEvent},
	}

	eventBody := EventV2Body{
		UserCreatedEvent: usersCreatedEvent,
		Challenge:        "1212",
		Type:             "url_verification",
	}

	body1, _ := json.Marshal(eventBody)

	return body1
}

func mockAppTicketEvent() []byte {

	body := "{\"ts\":\"\",\"uuid\":\"\",\"token\":\"1212121212\",\"type\":\"\",\"event\":{\"app_id\":\"jiaduoappId\",\"type\":\"app_ticket\",\"app_ticket\":\"AppTicketvalue\"}}"
	return []byte(body)
}

func main() {

	//mock body
	encryptedKey := "1212121212"
	body := mockEncryptedBody(encryptedKey)
	//body := mockEvent()

	// 创建http req
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9999/webhook/event", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 计算签名
	var timestamp = "timestamp"
	var nonce = "nonce"
	var token = encryptedKey
	sourceSign := event.Signature(timestamp, nonce, token, string(body))

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
