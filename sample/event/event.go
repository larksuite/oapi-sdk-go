package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/httpserverext"
	"github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func main() {

	//1212121212
	handler := dispatcher.NewEventDispatcher("verificationToken", "eventEncryptKey").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
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
	})

	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 开发者启动服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
