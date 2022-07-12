package larkapplication

import "context"

type P1OrderPaidV6Handler struct {
	handler func(context.Context, *P1OrderPaidV6) error
}

func NewP1OrderPaidV6Handler(handler func(context.Context, *P1OrderPaidV6) error) *P1OrderPaidV6Handler {
	h := &P1OrderPaidV6Handler{handler: handler}
	return h
}

func (h *P1OrderPaidV6Handler) Event() interface{} {
	return &P1OrderPaidV6{}
}

func (h *P1OrderPaidV6Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1OrderPaidV6))
}

type P1AppUninstalledV6Handler struct {
	handler func(context.Context, *P1AppUninstalledV6) error
}

func NewP1AppUninstalledV6Handler(handler func(context.Context, *P1AppUninstalledV6) error) *P1AppUninstalledV6Handler {
	h := &P1AppUninstalledV6Handler{handler: handler}
	return h
}

func (h *P1AppUninstalledV6Handler) Event() interface{} {
	return &P1AppUninstalledV6{}
}

func (h *P1AppUninstalledV6Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1AppUninstalledV6))
}

type P1AppStatusChangedV6Handler struct {
	handler func(context.Context, *P1AppStatusChangedV6) error
}

func NewP1AppStatusChangedV6Handler(handler func(context.Context, *P1AppStatusChangedV6) error) *P1AppStatusChangedV6Handler {
	h := &P1AppStatusChangedV6Handler{handler: handler}
	return h
}

func (h *P1AppStatusChangedV6Handler) Event() interface{} {
	return &P1AppStatusChangedV6{}
}

func (h *P1AppStatusChangedV6Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1AppStatusChangedV6))
}

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
