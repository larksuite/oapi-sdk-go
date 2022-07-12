package larkapplication

import "context"

type P1AppOpenV6Handler struct {
	handler func(context.Context, *P1AppOpenV6) error
}

func NewP1AppOpenV6Handler(handler func(context.Context, *P1AppOpenV6) error) *P1AppOpenV6Handler {
	h := &P1AppOpenV6Handler{handler: handler}
	return h
}

func (h *P1AppOpenV6Handler) Event() interface{} {
	return &P1AppOpenV6{}
}

func (h *P1AppOpenV6Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1AppOpenV6))
}
