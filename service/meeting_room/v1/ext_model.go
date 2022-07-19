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

package larkmeeting_room

import "github.com/larksuite/oapi-sdk-go/v3/event"

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
