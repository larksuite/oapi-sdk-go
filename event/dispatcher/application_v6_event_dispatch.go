// Package dispatcher code generated by oapi sdk gen
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

package dispatcher

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/v3/service/application/v6"
)

// 应用创建
//
// - 当企业内有新的应用被创建时推送此事件
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/application/events/created
func (dispatcher *EventDispatcher) OnP2ApplicationCreatedV6(handler func(ctx context.Context, event *larkapplication.P2ApplicationCreatedV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.created_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.created_v6")
	}
	dispatcher.eventType2EventHandler["application.application.created_v6"] = larkapplication.NewP2ApplicationCreatedV6Handler(handler)
	return dispatcher
}

// 应用审核
//
// - 通过订阅该事件，可接收应用审核（通过 / 拒绝）事件
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/application-app_version/events/audit
func (dispatcher *EventDispatcher) OnP2ApplicationAppVersionAuditV6(handler func(ctx context.Context, event *larkapplication.P2ApplicationAppVersionAuditV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.app_version.audit_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.app_version.audit_v6")
	}
	dispatcher.eventType2EventHandler["application.application.app_version.audit_v6"] = larkapplication.NewP2ApplicationAppVersionAuditV6Handler(handler)
	return dispatcher
}

// 申请发布应用
//
// - 通过订阅该事件，可接收应用提交发布申请事件
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/application-app_version/events/publish_apply
func (dispatcher *EventDispatcher) OnP2ApplicationAppVersionPublishApplyV6(handler func(ctx context.Context, event *larkapplication.P2ApplicationAppVersionPublishApplyV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.app_version.publish_apply_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.app_version.publish_apply_v6")
	}
	dispatcher.eventType2EventHandler["application.application.app_version.publish_apply_v6"] = larkapplication.NewP2ApplicationAppVersionPublishApplyV6Handler(handler)
	return dispatcher
}

// 撤回应用发布申请
//
// - 通过订阅该事件，可接收应用撤回发布申请事件
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/application-app_version/events/publish_revoke
func (dispatcher *EventDispatcher) OnP2ApplicationAppVersionPublishRevokeV6(handler func(ctx context.Context, event *larkapplication.P2ApplicationAppVersionPublishRevokeV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.app_version.publish_revoke_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.app_version.publish_revoke_v6")
	}
	dispatcher.eventType2EventHandler["application.application.app_version.publish_revoke_v6"] = larkapplication.NewP2ApplicationAppVersionPublishRevokeV6Handler(handler)
	return dispatcher
}

// 新增应用反馈
//
// - 当应用收到新反馈时，触发该事件
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/application-feedback/events/created
func (dispatcher *EventDispatcher) OnP2ApplicationFeedbackCreatedV6(handler func(ctx context.Context, event *larkapplication.P2ApplicationFeedbackCreatedV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.feedback.created_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.feedback.created_v6")
	}
	dispatcher.eventType2EventHandler["application.application.feedback.created_v6"] = larkapplication.NewP2ApplicationFeedbackCreatedV6Handler(handler)
	return dispatcher
}

// 反馈更新
//
// - 当反馈的处理状态被更新时，触发该事件
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/application-feedback/events/updated
func (dispatcher *EventDispatcher) OnP2ApplicationFeedbackUpdatedV6(handler func(ctx context.Context, event *larkapplication.P2ApplicationFeedbackUpdatedV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.feedback.updated_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.feedback.updated_v6")
	}
	dispatcher.eventType2EventHandler["application.application.feedback.updated_v6"] = larkapplication.NewP2ApplicationFeedbackUpdatedV6Handler(handler)
	return dispatcher
}

// -
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/event/app-availability-scope-extended
func (dispatcher *EventDispatcher) OnP2ApplicationVisibilityAddedV6(handler func(ctx context.Context, event *larkapplication.P2ApplicationVisibilityAddedV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.visibility.added_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.visibility.added_v6")
	}
	dispatcher.eventType2EventHandler["application.application.visibility.added_v6"] = larkapplication.NewP2ApplicationVisibilityAddedV6Handler(handler)
	return dispatcher
}

// 机器人自定义菜单
//
// - 当用户点击类型为事件的机器人菜单时触发
//
// - 事件描述文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/application-v6/bot/events/menu
func (dispatcher *EventDispatcher) OnP2BotMenuV6(handler func(ctx context.Context, event *larkapplication.P2BotMenuV6) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.bot.menu_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.bot.menu_v6")
	}
	dispatcher.eventType2EventHandler["application.bot.menu_v6"] = larkapplication.NewP2BotMenuV6Handler(handler)
	return dispatcher
}