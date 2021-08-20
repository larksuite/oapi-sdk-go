package v1

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/v2"
)

type messageReceiveEventHandler struct {
	handler func(context.Context, *lark.RawRequest, *MessageReceiveEvent) error
}

func (h *messageReceiveEventHandler) Event() interface{} {
	return &MessageReceiveEvent{}
}

func (h *messageReceiveEventHandler) Handle(ctx context.Context, req *lark.RawRequest, event interface{}) error {
	return h.handler(ctx, req, event.(*MessageReceiveEvent))
}

func (messageService *MessageService) ReceiveEventHandler(handler func(ctx context.Context, req *lark.RawRequest, event *MessageReceiveEvent) error) {
	lark.EventHandler(messageService.service.app, "im.message.receive_v1", &messageReceiveEventHandler{handler: handler})
}
