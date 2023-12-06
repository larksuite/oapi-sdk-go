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
	"context"
	"fmt"
	"net/http"

	larkhelpdesk "github.com/larksuite/oapi-sdk-go/v3/service/helpdesk/v1"

	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/service/application/v6"
	"github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	"github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/larksuite/oapi-sdk-go/v3/service/meeting_room/v1"
)

func main() {
	//1212121212
	handler := dispatcher.NewEventDispatcher("", "1212121212")
	handler.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1MessageReadV1(func(ctx context.Context, event *larkim.P1MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1MessageReceiveV1(func(ctx context.Context, event *larkim.P1MessageReceiveV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1UserStatusChangedV3(func(ctx context.Context, event *larkcontact.P1UserStatusChangedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1UserChangedV3(func(ctx context.Context, event *larkcontact.P1UserChangedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1DepartmentChangedV3(func(ctx context.Context, event *larkcontact.P1DepartmentChangedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1ContactScopeChangedV3(func(ctx context.Context, event *larkcontact.P1ContactScopeChangedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1AddBotV1(func(ctx context.Context, event *larkim.P1AddBotV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1RemoveAddBotV1(func(ctx context.Context, event *larkim.P1RemoveBotV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1UserInOutChatV1(func(ctx context.Context, event *larkim.P1UserInOutChatV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1ChatDisbandV1(func(ctx context.Context, event *larkim.P1ChatDisbandV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1GroupSettingUpdatedV1(func(ctx context.Context, event *larkim.P1GroupSettingUpdatedV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1AppOpenV6(func(ctx context.Context, event *larkapplication.P1AppOpenV6) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1P2PChatCreatedV1(func(ctx context.Context, event *larkim.P1P2PChatCreatedV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1ThirdPartyMeetingRoomChangedV1(func(ctx context.Context, event *larkmeeting_room.P1ThirdPartyMeetingRoomChangedV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1LeaveApprovalV4(func(ctx context.Context, event *larkapproval.P1LeaveApprovalV4) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1WorkApprovalV4(func(ctx context.Context, event *larkapproval.P1WorkApprovalV4) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1ShiftApprovalV4(func(ctx context.Context, event *larkapproval.P1ShiftApprovalV4) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1RemedyApprovalV4(func(ctx context.Context, event *larkapproval.P1RemedyApprovalV4) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1TripApprovalV4(func(ctx context.Context, event *larkapproval.P1TripApprovalV4) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1OutApprovalV4(func(ctx context.Context, event *larkapproval.P1OutApprovalV4) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1AppStatusChangedV6(func(ctx context.Context, event *larkapplication.P1AppStatusChangedV6) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1AppUninstalledV6(func(ctx context.Context, event *larkapplication.P1AppUninstalledV6) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP1OrderPaidV6(func(ctx context.Context, event *larkapplication.P1OrderPaidV6) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnCustomizedEvent("custom_event_type", func(ctx context.Context, event *larkevent.EventReq) error {
		// 原生消息体
		fmt.Println(string(event.Body))
		fmt.Println(larkcore.Prettify(event.Header))
		fmt.Println(larkcore.Prettify(event.RequestURI))
		fmt.Println(event.RequestId())

		// 处理消息
		cipherEventJsonStr, err := handler.ParseReq(ctx, event)
		if err != nil {
			//  错误处理
			return err
		}

		plainEventJsonStr, err := handler.DecryptEvent(ctx, cipherEventJsonStr)
		if err != nil {
			//  错误处理
			return err
		}

		// 处理解密后的 消息体
		fmt.Println(plainEventJsonStr)

		return nil
	}).OnP2TicketMessageCreatedV1(func(ctx context.Context, event *larkhelpdesk.P2TicketMessageCreatedV1) error {
		fmt.Println(event)
		return nil
	}).OnP2TicketUpdatedV1(func(ctx context.Context, event *larkhelpdesk.P2TicketUpdatedV1) error {
		fmt.Println(event)
		return nil
	})

	// 注册 http 路由
	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler,
		larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动服务
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}
