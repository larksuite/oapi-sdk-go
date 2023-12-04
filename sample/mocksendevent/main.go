/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	larkapplication "github.com/larksuite/oapi-sdk-go/v3/service/application/v6"
	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	"github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkmeeting_room "github.com/larksuite/oapi-sdk-go/v3/service/meeting_room/v1"
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
		Type:            "url_verification1",
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

func mockhelpdeskTicketUpdatedV1() []byte {
	body := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"bb8d3850b92f2cac9ce68862ebd63eb5\",\"token\":\"OGfWLx8LMOJYkaz3p4TxDgwvIxZWEcS3\",\"create_time\":\"1666593118336\",\"event_type\":\"helpdesk.ticket.updated_v1\",\"tenant_key\":\"736588c9260f175d\",\"app_id\":\"cli_9f8f27b375e9100c\"},\"event\":{\"object\":{\"chat_id\":\"oc_b3829240889c07a32126f8ad05891497\",\"closed_at\":1666593118000,\"created_at\":1666592755000,\"guest\":{\"id\":{\"open_id\":\"ou_dee0f0f124eca432bb4d1787025b3c4b\",\"union_id\":\"on_e434dd79941dfee7deb6fe3e41a8dd18\",\"user_id\":\"5f4a6a3b\"},\"name\":\"-\"},\"helpdesk_id\":\"6868205580868141059\",\"solve\":1,\"stage\":1,\"status\":50,\"ticket_id\":\"7157961378652045340\",\"updated_at\":1666593118000},\"old_object\":{\"status\":1,\"updated_at\":1666592755000}}}"
	return []byte(body)
}

func mockhelpdeskTicketCreatedV1() []byte {
	body := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"7b378539bfbdd0a94a18af7a374cff47\",\"token\":\"OGfWLx8LMOJYkaz3p4TxDgwvIxZWEcS3\",\"create_time\":\"1666583338473\",\"event_type\":\"helpdesk.ticket_message.created_v1\",\"tenant_key\":\"736588c9260f175d\",\"app_id\":\"cli_9f8f27b375e9100c\"},\"event\":{\"chat_id\":\"oc_3f37a0031825d258e3f66d5a8b600c7d\",\"content\":{\"content\":\"**猜你想问**\\n--------\\n[安装使用哪些软件属于盗版软件？]\\n[腾讯相关软件下载（微信/腾讯会议/QQ/QQ浏览器等）]\\n[以上都不是，转人工服务]\\n\",\"msg_type\":\"text\"},\"event_id\":\"debbc874-1c84-b9bf-55cd-a8ab3dfb4e0c\",\"message_id\":\"om_4b95319c9e249e4e60800ec08f863bee\",\"msg_type\":\"text\",\"position\":222,\"sender_type\":1,\"text\":\"**猜你想问**\\n--------\\n[安装使用哪些软件属于盗版软件？]\\n[腾讯相关软件下载（微信/腾讯会议/QQ/QQ浏览器等）]\\n[以上都不是，转人工服务]\\n\",\"ticket\":{\"status\":1,\"ticket_id\":\"7157919461389959171\"},\"ticket_message_id\":\"7157920934710771713\"}}"
	return []byte(body)
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
			Ts:    "t",
			UUID:  "u",
			Token: "v",
			Type:  "type",
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
			Users: []*larkim.P1UserV1{{
				OpenId: "o1",
				UserId: "u1",
				Name:   "n1",
			}, {
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

func mockAppOpenEventV1() []byte {
	event := &larkapplication.P1AppOpenV6{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapplication.P1AppOpenV6Data{
			Type:      "app_open",
			AppID:     "xxx",
			TenantKey: "tk",
			Applicants: []*larkapplication.P1AppOpenApplicantV6{
				{OpenID: "o1"}, {OpenID: "o2"},
			},
			Installer:         &larkapplication.P1AppOpenInstallerV6{OpenID: "o1"},
			InstallerEmployee: &larkapplication.P1AppOpenInstallerEmployeeV6{OpenID: "o1"},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockP2PCreatedChatEventV1() []byte {
	event := &larkim.P1P2PChatCreatedV1{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkim.P1P2PChatCreatedV1Data{
			Type:      "p2p_chat_create",
			AppID:     "xxx",
			TenantKey: "tk",
			ChatID:    "chatid",
			Operator: &larkim.P1OperatorV1{
				OpenId: "oi",
				UserId: "ui",
			},
			User: &larkim.P1UserV1{
				OpenId: "oi",
				UserId: "ui",
				Name:   "jiaduo",
			},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockThirdMeetingRoomChangedEventV1() []byte {
	event := &larkmeeting_room.P1ThirdPartyMeetingRoomChangedV1{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkmeeting_room.P1ThirdPartyMeetingRoomChangedV1Data{
			Type:         "third_party_meeting_room_event_updated",
			AppID:        "xxx",
			TenantKey:    "tk",
			EventTime:    "1594979647635",
			Uid:          "bff6b51f",
			OriginalTime: 0,
			Start:        &larkmeeting_room.P1EventTimeV1{TimeStamp: "1553853600000"},
			End:          &larkmeeting_room.P1EventTimeV1{TimeStamp: "1553860800000"},
			MeetingRoom:  []*larkmeeting_room.P1MeetingRoomV1{{OpenId: "oi1"}, {OpenId: "oi2"}},
			Organizer: &larkmeeting_room.P1OrganizerV1{
				OpenId: "oi",
				UserId: "ui",
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

func mockLeaveApprovalEventV1() []byte {
	event := &larkapproval.P1LeaveApprovalV4{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapproval.P1LeaveApprovalV4Data{
			Type:                   "leave_approvalV2",
			AppID:                  "xxx",
			TenantKey:              "tk",
			InstanceCode:           "code",
			UserID:                 "userid",
			OpenID:                 "openid",
			OriginInstanceCode:     "origincode",
			StartTime:              1564590532,
			EndTime:                1564590533,
			LeaveFeedingArriveLate: 1,
			LeaveFeedingLeaveEarly: 2,
			LeaveFeedingRestDaily:  3,
			LeaveName:              "JIADUO",
			LeaveUnit:              "day",
			LeaveStartTime:         "2019-10-01 00:00:00",
			LeaveEndTime:           "2019-10-01 00:00:00",
			LeaveDetail:            []string{"2019-10-01 00:00:00", "2019-10-02 00:00:00"},
			LeaveRange:             []string{"2019-10-01 00:00:00", "2019-10-02 00:00:00"},
			LeaveInterval:          86400,
			LeaveReason:            "abc",
			I18nResources: []*larkapproval.P1LeaveApprovalI18nResourceV4{
				{Locale: "en_us",
					IsDefault: true,
					Texts: map[string]string{
						"@i18n@123456": "Holiday",
					},
				},
			},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockWorkApprovalEventV1() []byte {
	event := &larkapproval.P1WorkApprovalV4{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapproval.P1WorkApprovalV4Data{
			Type:          "work_approval",
			AppID:         "xxx",
			TenantKey:     "tk",
			InstanceCode:  "code",
			OpenID:        "openid",
			StartTime:     1564590532,
			EndTime:       1564590533,
			EmployeeID:    "id",
			WorkType:      "ss",
			WorkStartTime: "2018-12-01 12:00:00",
			WorkEndTime:   "2018-12-03 12:00:00",
			WorkInterval:  1000,
			WorkReason:    "reason",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockShiftApprovalEventV1() []byte {
	event := &larkapproval.P1ShiftApprovalV4{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapproval.P1ShiftApprovalV4Data{
			Type:         "shift_approval",
			AppID:        "xxx",
			TenantKey:    "tk",
			InstanceCode: "code",
			OpenID:       "openid",
			StartTime:    1564590532,
			EndTime:      1564590533,
			EmployeeID:   "id",
			ShiftTime:    "2018-12-01 12:00:00",
			ReturnTime:   "2018-12-01 12:00:00",
			ShiftReason:  "reason",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockRemedyApprovalEventV1() []byte {
	event := &larkapproval.P1RemedyApprovalV4{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapproval.P1RemedyApprovalV4Data{
			Type:         "remedy_approval",
			AppID:        "xxx",
			TenantKey:    "tk",
			InstanceCode: "code",
			OpenID:       "openid",
			StartTime:    1564590532,
			EndTime:      1564590533,
			EmployeeID:   "id",
			RemedyReason: "reason",
			RemedyTime:   "0",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockTripApprovalEventV1() []byte {
	event := &larkapproval.P1TripApprovalV4{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapproval.P1TripApprovalV4Data{
			Type:         "trip_approval",
			AppID:        "xxx",
			TenantKey:    "tk",
			InstanceCode: "code",
			OpenID:       "openid",
			StartTime:    1564590532,
			EndTime:      1564590533,
			EmployeeID:   "id",
			TripInterval: 12121,
			TripReason:   "reason",
			TripPeers:    []string{"a", "b"},
			Schedules: []*larkapproval.P1TripApprovalScheduleV4{
				{
					TripStartTime:  "2018-12-01 12:00:00",
					TripEndTime:    "2018-12-02 12:00:00",
					TripInterval:   3600,
					Departure:      "xxx",
					Destination:    "x",
					Transportation: "xxxx",
					TripType:       "单程",
					Remark:         "备注",
				},
			},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockOutApprovalEventV1() []byte {
	event := &larkapproval.P1OutApprovalV4{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapproval.P1OutApprovalV4Data{
			Type:         "out_approval",
			AppID:        "xxx",
			TenantKey:    "tk",
			InstanceCode: "code",
			OpenID:       "openid",
			StartTime:    1564590532,
			EndTime:      1564590533,
			I18nResources: []*larkapproval.P1OutApprovalI18nResourceV4{{
				Locale:    "en_us",
				IsDefault: true,
				Texts:     map[string]string{"k1": "v1", "k2": "v2"},
			}},
			OutImage:     "image",
			OutInterval:  1000,
			OutName:      "name",
			OutReason:    "事由",
			OutStartTime: "2020-05-15 15:00:00",
			OutEndTime:   "2020-05-16 15:00:00",
			OutUnit:      "HOUR",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockAppStatusChangedEventV1() []byte {
	event := &larkapplication.P1AppStatusChangedV6{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapplication.P1AppStatusChangedV6Data{
			Type:      "app_status_change",
			AppID:     "xxx",
			TenantKey: "tk",
			Status:    "start_by_tenant",
			Operator: &larkapplication.P1AppStatusChangeOperatorV6{
				OpenID:  "o1",
				UserID:  "ui",
				UnionId: "ui",
			},
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockAppUninstalledEventV1() []byte {
	event := &larkapplication.P1AppUninstalledV6{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "",
		},
		Event: &larkapplication.P1AppUninstalledV6Data{
			Type:      "app_uninstalled",
			AppID:     "xxx",
			TenantKey: "tk",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func mockAppOrderPaidEventV1() []byte {
	event := &larkapplication.P1OrderPaidV6{
		EventBase: &larkevent.EventBase{
			Ts:    "ts",
			UUID:  "uid",
			Token: "v",
			Type:  "type",
		},
		Event: &larkapplication.P1OrderPaidV6Data{
			Type:          "order_paid",
			AppID:         "xxx",
			TenantKey:     "tk",
			OrderID:       "54323223",
			PricePlanID:   "price_12121",
			PricePlanType: "per_seat_per_month",
			Seats:         20,
			BuyType:       "buy",
			BuyCount:      1,
			SrcOrderID:    "23233",
			OrderPayPrice: 10000,
			CreateTime:    "time",
			PayTime:       "paytime",
		},
	}
	body1, _ := json.Marshal(event)
	return body1
}

func main() {

	//mock body
	encryptedKey := "1212121212"
	//body := mockMessageReceiveEventV1()
	//body := mockAppTicketEvent()
	//body := mockMessageReceiveEventV1()
	//body := mockAppTicketEvent()
	body := mockEncryptedBody(encryptedKey)

	// 创建http req
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:7777/webhook/event", bytes.NewBuffer(body))
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
