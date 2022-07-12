package larkmeeting_room

import "github.com/larksuite/oapi-sdk-go/event"

type P1ThirdPartyMeetingRoomChangedV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1ThirdPartyMeetingRoomChangedV1Data `json:"event"`
}

func (m *P1ThirdPartyMeetingRoomChangedV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1EventTimeV1 struct {
	TimeStamp string `json:"time_stamp,omitempty"`
}

type P1MeetingRoomV1 struct {
	OpenId string `json:"open_id,omitempty"`
}

type P1OrganizerV1 struct {
	OpenId string `json:"open_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

type P1ThirdPartyMeetingRoomChangedV1Data struct {
	AppID        string             `json:"app_id,omitempty"`
	TenantKey    string             `json:"tenant_key,omitempty"`
	Type         string             `json:"type,omitempty"`
	EventTime    string             `json:"event_time,omitempty"`
	Uid          string             `json:"uid,omitempty"`
	OriginalTime int                `json:"original_time,omitempty"`
	Start        *P1EventTimeV1     `json:"start,omitempty"`
	End          *P1EventTimeV1     `json:"end,omitempty"`
	MeetingRoom  []*P1MeetingRoomV1 `json:"meeting_rooms,omitempty"`
	Organizer    *P1OrganizerV1     `json:"organizer,omitempty"`
}
