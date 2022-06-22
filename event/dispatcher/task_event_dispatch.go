// Package dispatcher code generated by oapi sdk gen
package dispatcher

import (
	"context"

	"github.com/larksuite/oapi-sdk-go/service/task/v1"
)

func (dispatcher *EventDispatcher) OnTaskUpdateTenantV1(handler func(ctx context.Context, event *task.TaskUpdateTenantEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["task.task.update_tenant_v1"]
	if existed {
		panic("event: multiple handler registrations for " + "task.task.update_tenant_v1")
	}
	dispatcher.eventType2EventHandler["task.task.update_tenant_v1"] = task.NewTaskUpdateTenantEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnTaskUpdatedV1(handler func(ctx context.Context, event *task.TaskUpdatedEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["task.task.updated_v1"]
	if existed {
		panic("event: multiple handler registrations for " + "task.task.updated_v1")
	}
	dispatcher.eventType2EventHandler["task.task.updated_v1"] = task.NewTaskUpdatedEventHandler(handler)
	return dispatcher
}
func (dispatcher *EventDispatcher) OnTaskCommentUpdatedV1(handler func(ctx context.Context, event *task.TaskCommentUpdatedEvent) error) *EventDispatcher {
	_, existed := dispatcher.eventType2EventHandler["task.task.comment.updated_v1"]
	if existed {
		panic("event: multiple handler registrations for " + "task.task.comment.updated_v1")
	}
	dispatcher.eventType2EventHandler["task.task.comment.updated_v1"] = task.NewTaskCommentUpdatedEventHandler(handler)
	return dispatcher
}