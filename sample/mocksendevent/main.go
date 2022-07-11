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
	larkim "github.com/larksuite/oapi-sdk-go/service/im/v1"
)

type EventV2Body struct {
	larkcontact.P2UserCreatedV3        // 什么类型消息就替换为什么
	Challenge                   string `json:"challenge"`
	Type                        string `json:"type"`
	ts                          string `json:"ts"`
	uuid                        string `json:"uuid"`
	token                       string `json:"token"`
}

func mockEncryptedBody(encrypteKey string) []byte {

	userEvent := &larkcontact.UserEvent{
		OpenId:     larkcore.StringPtr("ou_7dab8a3d3cdcc9da365777c7ad535d62"),
		UnionId:    larkcore.StringPtr("on_576833b917gda3d939b9a3c2d53e72c8"),
		UserId:     larkcore.StringPtr("e33ggbyz"),
		Name:       larkcore.StringPtr("张三"),
		EmployeeNo: larkcore.StringPtr("employee_no"),
	}

	usersCreatedEvent := larkcontact.P2UserCreatedV3{
		EventV2Base: &larkevent.EventV2Base{
			Schema: "2.0",
			Header: &larkevent.EventHeader{
				EventID:    "f7984f25108f8137722bb63cee927e66",
				EventType:  "contact.user.created_v3",
				AppID:      "cli_xxxxxxxx",
				TenantKey:  "xxxxxxx",
				CreateTime: "1603977298000000",
				Token:      "v",
			},
		},
		Event: &larkcontact.P2UserCreatedV3Data{Object: userEvent},
	}

	eventBody := EventV2Body{
		P2UserCreatedV3: usersCreatedEvent,
		Challenge:       "1212",
		Type:            "url_verification",
	}

	en, _ := larkcore.EncryptedEventMsg(context.Background(), eventBody, encrypteKey)
	fmt.Println(encrypteKey)

	encrypt := larkevent.EventEncryptMsg{Encrypt: en}
	body1, _ := json.Marshal(encrypt)

	return body1
}

func mockEvent() []byte {
	userEvent := &larkcontact.UserEvent{
		OpenId:     larkcore.StringPtr("ou_7dab8a3d3cdcc9da365777c7ad535d62"),
		UnionId:    larkcore.StringPtr("on_576833b917gda3d939b9a3c2d53e72c8"),
		UserId:     larkcore.StringPtr("e33ggbyz"),
		Name:       larkcore.StringPtr("张三"),
		EmployeeNo: larkcore.StringPtr("employee_no"),
	}

	usersCreatedEvent := larkcontact.P2UserCreatedV3{
		EventV2Base: &larkevent.EventV2Base{
			Schema: "2.0",
			Header: &larkevent.EventHeader{
				EventID:    "f7984f25108f8137722bb63cee927e66",
				EventType:  "contact.user.created_v3",
				AppID:      "cli_xxxxxxxx",
				TenantKey:  "xxxxxxx",
				CreateTime: "1603977298000000",
				Token:      "v",
			},
		},
		Event: &larkcontact.P2UserCreatedV3Data{Object: userEvent},
	}

	eventBody := EventV2Body{
		P2UserCreatedV3: usersCreatedEvent,
		Challenge:       "1212",
		Type:            "url_verification1",
	}

	body1, _ := json.Marshal(eventBody)

	return body1
}

type EventFuzzy struct {
	Encrypt   string       `json:"encrypt"`
	Schema    string       `json:"schema"`
	Token     string       `json:"token"`
	uuid      string       `json:"uuid"`
	Type      string       `json:"type"`
	Challenge string       `json:"challenge"`
	Event     *EventV1Body `json:"event"`
}

type EventV1Body struct {
	*larkcontact.UserEvent        // 什么类型消息就替换为什么
	Type                   string `json:"type"`
}

type MessageReadEventV1Body struct {
	*larkim.P1MessageReadV1 // 什么类型消息就替换为什么
}

func mockEventV1() []byte {
	userEvent := &larkcontact.UserEvent{
		OpenId:     larkcore.StringPtr("ou_7dab8a3d3cdcc9da365777c7ad535d62"),
		UnionId:    larkcore.StringPtr("on_576833b917gda3d939b9a3c2d53e72c8"),
		UserId:     larkcore.StringPtr("e33ggbyz"),
		Name:       larkcore.StringPtr("张三"),
		EmployeeNo: larkcore.StringPtr("employee_no"),
	}

	body := EventFuzzy{
		Challenge: "1212",
		Type:      "url_verification1",
		uuid:      "sssss",
		Token:     "1",
		Event: &EventV1Body{
			UserEvent: userEvent,
			Type:      "user_add",
		},
	}

	body1, _ := json.Marshal(body)

	return body1
}

func mockMessageReadEventV1() []byte {
	event := &larkim.P1MessageReadV1{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1MessageReadV1Data{
			MessageIdList: []string{"ss", "dd"},
			AppID:         "appid",
			OpenAppID:     "openapiid",
			OpenID:        "openid",
			TenantKey:     "tenkey",
			Type:          "message_read",
		},
	}

	body1, _ := json.Marshal(event)

	return body1
}

func mockMessageReceiveEventV1() []byte {
	event := &larkim.P1MessageReceiveV1{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1MessageReceiveV1Data{
			Type:             "message",
			AppID:            "appid",
			TenantKey:        "tenantkey",
			RootID:           "1212",
			ParentID:         "11",
			OpenChatID:       "123",
			ChatType:         "public",
			MsgType:          "text",
			OpenID:           "1221",
			EmployeeID:       "sdsd",
			UnionID:          "sdsd",
			OpenMessageID:    "2wwd",
			IsMention:        false,
			Text:             "hello jiaduo",
			TextWithoutAtBot: "",
			Title:            "title",
			ImageKeys:        nil,
			ImageKey:         "",
			FileKey:          "",
		},
	}

	body1, _ := json.Marshal(event)

	return body1
}

func mockUserStatusChangedEventV1() []byte {
	event := &larkcontact.P1UserStatusChangedV3{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkcontact.P1UserStatusChangedV3Data{
			Type:       "user_status_change",
			AppID:      "cli_xxx",
			TenantKey:  "xx",
			OpenID:     "xx",
			EmployeeId: "xx",
			UnionId:    "xx",
			BeforeStatus: &larkcontact.P1UserStatusV3{
				IsActive:   false,
				IsFrozen:   false,
				IsResigned: false,
			},
			CurrentStatus: &larkcontact.P1UserStatusV3{
				IsActive:   true,
				IsFrozen:   false,
				IsResigned: false,
			},
			ChangeTime: "2020-02-21 16:28:48",
		},
	}

	body1, _ := json.Marshal(event)

	return body1
}

func mockUserChangedEventV1() []byte {
	event := &larkcontact.P1UserChangedV3{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkcontact.P1UserChangedV3Data{
			Type:       "user_leave",
			AppID:      "cli_xxx",
			TenantKey:  "xx",
			OpenID:     "xx",
			EmployeeId: "xx",
			UnionId:    "xx",
		},
	}

	body1, _ := json.Marshal(event)

	return body1
}

func mockDeptChangedEventV1() []byte {
	event := &larkcontact.P1DepartmentChangedV3{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkcontact.P1DepartmentChangedV3Data{
			Type:             "dept_update",
			AppID:            "cli_xxx",
			TenantKey:        "xx",
			OpenID:           "xx",
			OpenDepartmentId: "sssss",
		},
	}

	body1, _ := json.Marshal(event)

	return body1
}

func mockContactScopeChangedEventV1() []byte {
	event := &larkcontact.P1ContactScopeChangedV3{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkcontact.P1ContactScopeChangedV3Data{
			Type:      "contact_scope_change",
			AppID:     "cli_xxx",
			TenantKey: "xx",
		},
	}

	body1, _ := json.Marshal(event)

	return body1
}

func mockAddBotEventV1() []byte {
	event := &larkim.P1AddBotV1{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1AddBotV1Data{
			Type:  "add_bot",
			AppID: "xxx",
			ChatI18nNames: &larkim.ChatI18nNames{
				EnUs: "zlx",
				ZhCn: "加多",
			},
			ChatName:            "告警",
			ChatOwnerEmployeeID: "121212",
			ChatOwnerName:       "name",
			ChatOwnerOpenID:     "opneid",
			OpenChatID:          "chatid",
			OperatorEmployeeID:  "eid",
			OperatorName:        "on",
			OperatorOpenID:      "oi",
			OwnerIsBot:          false,
			TenantKey:           "tk",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockRemoveBotEventV1() []byte {
	event := &larkim.P1RemoveBotV1{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1RemoveBotV1Data{
			Type:  "remove_bot",
			AppID: "xxx",
			ChatI18nNames: &larkim.ChatI18nNames{
				EnUs: "zlx",
				ZhCn: "加多",
			},
			ChatName:            "告警",
			ChatOwnerEmployeeID: "121212",
			ChatOwnerName:       "name",
			ChatOwnerOpenID:     "opneid",
			OpenChatID:          "chatid",
			OperatorEmployeeID:  "eid",
			OperatorName:        "on",
			OperatorOpenID:      "oi",
			OwnerIsBot:          false,
			TenantKey:           "tk",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockUserInOutChatEventV1() []byte {
	event := &larkim.P1UserInOutChatV1{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1UserInOutChatV1Data{
			Type:      "revoke_add_user_from_chat",
			AppID:     "xxx",
			TenantKey: "tk",
			ChatId:    "chatid",
			Operator: &larkim.P1OperatorV1{
				OpenId: "openid",
				UserId: "userid",
			},
			Users: []*larkim.P1UserV1{&larkim.P1UserV1{
				OpenId: "o1",
				UserId: "u1",
				Name:   "n1",
			}, &larkim.P1UserV1{
				OpenId: "o2",
				UserId: "u2",
				Name:   "n2",
			}},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockDisbandChatEventV1() []byte {
	event := &larkim.P1ChatDisbandV1{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1ChatDisbandV1Data{
			Type:      "chat_disband",
			AppID:     "xxx",
			TenantKey: "tk",
			ChatId:    "chatid",
			Operator: &larkim.P1OperatorV1{
				OpenId: "openid",
				UserId: "userid",
			},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockGroupSettingUpdatedEventV1() []byte {
	event := &larkim.P1GroupSettingUpdatedV1{
		EventBase: &larkevent.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1GroupSettingUpdatedV1Data{
			Type:      "group_setting_update",
			AppID:     "xxx",
			TenantKey: "tk",
			ChatId:    "chatid",
			Operator: &larkim.P1OperatorV1{
				OpenId: "openid",
				UserId: "userid",
			},
			BeforeChange: &larkim.P1GroupSettingChangeV1{
				OwnerOpenId:         "ooi",
				OwnerUserId:         "oui",
				AddMemberPermission: "amp",
				MessageNotification: false,
			},
			AfterChange: &larkim.P1GroupSettingChangeV1{
				OwnerOpenId:         "ooi",
				OwnerUserId:         "oui",
				AddMemberPermission: "amp",
				MessageNotification: true,
			},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockAppTicketEvent() []byte {
	body := "{\"ts\":\"\",\"uuid\":\"\",\"token\":\"1212121212\",\"type\":\"\",\"event\":{\"app_id\":\"jiaduoappId\",\"type\":\"app_ticket\",\"app_ticket\":\"AppTicketvalue\"}}"
	return []byte(body)
}

func main() {

	//mock body
	encryptedKey := "1212121212"
	//body := mockEncryptedBody(encryptedKey)
	//body := mockEvent()
	//body := mockMessageReceiveEventV1()
	//body := mockAppTicketEvent()
	body := mockEvent()

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
	sourceSign := larkevent.Signature(timestamp, nonce, token, string(body))

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
