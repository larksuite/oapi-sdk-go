package larkapproval

import "context"

type P1LeaveApprovalV4Handler struct {
	handler func(context.Context, *P1LeaveApprovalV4) error
}

func NewP1LeaveApprovalV4Handler(handler func(context.Context, *P1LeaveApprovalV4) error) *P1LeaveApprovalV4Handler {
	h := &P1LeaveApprovalV4Handler{handler: handler}
	return h
}

func (h *P1LeaveApprovalV4Handler) Event() interface{} {
	return &P1LeaveApprovalV4{}
}

func (h *P1LeaveApprovalV4Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1LeaveApprovalV4))
}

type P1WorkApprovalV4Handler struct {
	handler func(context.Context, *P1WorkApprovalV4) error
}

func NewP1WorkApprovalV4Handler(handler func(context.Context, *P1WorkApprovalV4) error) *P1WorkApprovalV4Handler {
	h := &P1WorkApprovalV4Handler{handler: handler}
	return h
}

func (h *P1WorkApprovalV4Handler) Event() interface{} {
	return &P1WorkApprovalV4{}
}

func (h *P1WorkApprovalV4Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1WorkApprovalV4))
}

type P1ShiftApprovalV4Handler struct {
	handler func(context.Context, *P1ShiftApprovalV4) error
}

func NewP1ShiftApprovalV4Handler(handler func(context.Context, *P1ShiftApprovalV4) error) *P1ShiftApprovalV4Handler {
	h := &P1ShiftApprovalV4Handler{handler: handler}
	return h
}

func (h *P1ShiftApprovalV4Handler) Event() interface{} {
	return &P1ShiftApprovalV4{}
}

func (h *P1ShiftApprovalV4Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1ShiftApprovalV4))
}

type P1RemedyApprovalV4Handler struct {
	handler func(context.Context, *P1RemedyApprovalV4) error
}

func NewP1RemedyApprovalV4Handler(handler func(context.Context, *P1RemedyApprovalV4) error) *P1RemedyApprovalV4Handler {
	h := &P1RemedyApprovalV4Handler{handler: handler}
	return h
}

func (h *P1RemedyApprovalV4Handler) Event() interface{} {
	return &P1RemedyApprovalV4{}
}

func (h *P1RemedyApprovalV4Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1RemedyApprovalV4))
}

type P1TripApprovalV4Handler struct {
	handler func(context.Context, *P1TripApprovalV4) error
}

func NewP1TripApprovalV4Handler(handler func(context.Context, *P1TripApprovalV4) error) *P1TripApprovalV4Handler {
	h := &P1TripApprovalV4Handler{handler: handler}
	return h
}

func (h *P1TripApprovalV4Handler) Event() interface{} {
	return &P1TripApprovalV4{}
}

func (h *P1TripApprovalV4Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1TripApprovalV4))
}

type P1OutApprovalV4Handler struct {
	handler func(context.Context, *P1OutApprovalV4) error
}

func NewP1OutApprovalV4Handler(handler func(context.Context, *P1OutApprovalV4) error) *P1OutApprovalV4Handler {
	h := &P1OutApprovalV4Handler{handler: handler}
	return h
}

func (h *P1OutApprovalV4Handler) Event() interface{} {
	return &P1OutApprovalV4{}
}

func (h *P1OutApprovalV4Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1OutApprovalV4))
}
