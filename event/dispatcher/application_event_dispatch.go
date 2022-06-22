// Package dispatcher code generated by oapi sdk gen
package dispatcher

import (
	"context"

	"github.com/larksuite/oapi-sdk-go/service/application/v6"
)

func (dispatcher *EventDispatcher) OnApplicationCreatedV6(handler func(ctx context.Context, event *application.ApplicationCreatedEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.created_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.created_v6")
	}
	dispatcher.eventType2EventHandler["application.application.created_v6"] = application.NewApplicationCreatedEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnApplicationAppVersionAuditV6(handler func(ctx context.Context, event *application.ApplicationAppVersionAuditEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.app_version.audit_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.app_version.audit_v6")
	}
	dispatcher.eventType2EventHandler["application.application.app_version.audit_v6"] = application.NewApplicationAppVersionAuditEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnApplicationAppVersionPublishApplyV6(handler func(ctx context.Context, event *application.ApplicationAppVersionPublishApplyEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.app_version.publish_apply_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.app_version.publish_apply_v6")
	}
	dispatcher.eventType2EventHandler["application.application.app_version.publish_apply_v6"] = application.NewApplicationAppVersionPublishApplyEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnApplicationAppVersionPublishRevokeV6(handler func(ctx context.Context, event *application.ApplicationAppVersionPublishRevokeEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.app_version.publish_revoke_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.app_version.publish_revoke_v6")
	}
	dispatcher.eventType2EventHandler["application.application.app_version.publish_revoke_v6"] = application.NewApplicationAppVersionPublishRevokeEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnApplicationFeedbackCreatedV6(handler func(ctx context.Context, event *application.ApplicationFeedbackCreatedEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.feedback.created_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.feedback.created_v6")
	}
	dispatcher.eventType2EventHandler["application.application.feedback.created_v6"] = application.NewApplicationFeedbackCreatedEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnApplicationFeedbackUpdatedV6(handler func(ctx context.Context, event *application.ApplicationFeedbackUpdatedEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.feedback.updated_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.feedback.updated_v6")
	}
	dispatcher.eventType2EventHandler["application.application.feedback.updated_v6"] = application.NewApplicationFeedbackUpdatedEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnApplicationVisibilityAddedV6(handler func(ctx context.Context, event *application.ApplicationVisibilityAddedEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["application.application.visibility.added_v6"]
	if existed {
		panic("event: multiple handler registrations for " + "application.application.visibility.added_v6")
	}
	dispatcher.eventType2EventHandler["application.application.visibility.added_v6"] = application.NewApplicationVisibilityAddedEventHandler(handler)
	return dispatcher
}