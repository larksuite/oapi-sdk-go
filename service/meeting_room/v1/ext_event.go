package larkmeeting_room

import "context"

type P1ThirdPartyMeetingRoomChangedV1Handler struct {
	handler func(context.Context, *P1ThirdPartyMeetingRoomChangedV1) error
}

func NewP1ThirdPartyMeetingRoomChangedV1Handler(handler func(context.Context, *P1ThirdPartyMeetingRoomChangedV1) error) *P1ThirdPartyMeetingRoomChangedV1Handler {
	h := &P1ThirdPartyMeetingRoomChangedV1Handler{handler: handler}
	return h
}

func (h *P1ThirdPartyMeetingRoomChangedV1Handler) Event() interface{} {
	return &P1ThirdPartyMeetingRoomChangedV1{}
}

func (h *P1ThirdPartyMeetingRoomChangedV1Handler) Handle(ctx context.Context, event interface{}) error {
	return h.handler(ctx, event.(*P1ThirdPartyMeetingRoomChangedV1))
}
