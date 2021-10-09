// Package calendar code generated by lark suite oapi sdk gen
package calendar

import (
	"context"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v2"
)

type CalendarService struct {
	Calendars                        *calendars
	CalendarAcls                     *calendarAcls
	CalendarEvents                   *calendarEvents
	CalendarEventAttendees           *calendarEventAttendees
	CalendarEventAttendeeChatMembers *calendarEventAttendeeChatMembers
	Freebusy                         *freebusy
	Settings                         *settings
	TimeoffEvents                    *timeoffEvents
}

func New(app *lark.App) *CalendarService {
	c := &CalendarService{}
	c.Calendars = &calendars{app: app}
	c.CalendarAcls = &calendarAcls{app: app}
	c.CalendarEvents = &calendarEvents{app: app}
	c.CalendarEventAttendees = &calendarEventAttendees{app: app}
	c.CalendarEventAttendeeChatMembers = &calendarEventAttendeeChatMembers{app: app}
	c.Freebusy = &freebusy{app: app}
	c.Settings = &settings{app: app}
	c.TimeoffEvents = &timeoffEvents{app: app}
	return c
}

type calendars struct {
	app *lark.App
}
type calendarAcls struct {
	app *lark.App
}
type calendarEvents struct {
	app *lark.App
}
type calendarEventAttendees struct {
	app *lark.App
}
type calendarEventAttendeeChatMembers struct {
	app *lark.App
}
type freebusy struct {
	app *lark.App
}
type settings struct {
	app *lark.App
}
type timeoffEvents struct {
	app *lark.App
}

func (c *calendars) Create(ctx context.Context, req *CalendarCreateReq, options ...lark.RequestOptionFunc) (*CalendarCreateResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarCreateResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) Patch(ctx context.Context, req *CalendarPatchReq, options ...lark.RequestOptionFunc) (*CalendarPatchResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPatch,
		"/open-apis/calendar/v4/calendars/:calendar_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarPatchResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) Delete(ctx context.Context, req *CalendarDeleteReq, options ...lark.RequestOptionFunc) (*CalendarDeleteResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodDelete,
		"/open-apis/calendar/v4/calendars/:calendar_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarDeleteResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) List(ctx context.Context, req *CalendarListReq, options ...lark.RequestOptionFunc) (*CalendarListResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodGet,
		"/open-apis/calendar/v4/calendars", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarListResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) Get(ctx context.Context, req *CalendarGetReq, options ...lark.RequestOptionFunc) (*CalendarGetResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodGet,
		"/open-apis/calendar/v4/calendars/:calendar_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarGetResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) Search(ctx context.Context, req *CalendarSearchReq, options ...lark.RequestOptionFunc) (*CalendarSearchResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/search", []lark.AccessTokenType{lark.AccessTokenTypeTenant}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarSearchResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) Unsubscribe(ctx context.Context, req *CalendarUnsubscribeReq, options ...lark.RequestOptionFunc) (*CalendarUnsubscribeResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/unsubscribe", []lark.AccessTokenType{lark.AccessTokenTypeUser, lark.AccessTokenTypeTenant}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarUnsubscribeResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) Subscribe(ctx context.Context, req *CalendarSubscribeReq, options ...lark.RequestOptionFunc) (*CalendarSubscribeResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/subscribe", []lark.AccessTokenType{lark.AccessTokenTypeUser, lark.AccessTokenTypeTenant}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarSubscribeResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendars) Subscription(ctx context.Context, options ...lark.RequestOptionFunc) (*CalendarSubscriptionResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/subscription", []lark.AccessTokenType{lark.AccessTokenTypeUser}, nil, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarSubscriptionResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarAcls) List(ctx context.Context, req *CalendarAclListReq, options ...lark.RequestOptionFunc) (*CalendarAclListResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodGet,
		"/open-apis/calendar/v4/calendars/:calendar_id/acls", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarAclListResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarAcls) Delete(ctx context.Context, req *CalendarAclDeleteReq, options ...lark.RequestOptionFunc) (*CalendarAclDeleteResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodDelete,
		"/open-apis/calendar/v4/calendars/:calendar_id/acls/:acl_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarAclDeleteResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarAcls) Create(ctx context.Context, req *CalendarAclCreateReq, options ...lark.RequestOptionFunc) (*CalendarAclCreateResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/acls", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarAclCreateResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarAcls) Subscription(ctx context.Context, req *CalendarAclSubscriptionReq, options ...lark.RequestOptionFunc) (*CalendarAclSubscriptionResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/acls/subscription", []lark.AccessTokenType{lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarAclSubscriptionResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEvents) Delete(ctx context.Context, req *CalendarEventDeleteReq, options ...lark.RequestOptionFunc) (*CalendarEventDeleteResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodDelete,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventDeleteResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEvents) Get(ctx context.Context, req *CalendarEventGetReq, options ...lark.RequestOptionFunc) (*CalendarEventGetResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodGet,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventGetResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEvents) Create(ctx context.Context, req *CalendarEventCreateReq, options ...lark.RequestOptionFunc) (*CalendarEventCreateResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/events", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventCreateResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEvents) List(ctx context.Context, req *CalendarEventListReq, options ...lark.RequestOptionFunc) (*CalendarEventListResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodGet,
		"/open-apis/calendar/v4/calendars/:calendar_id/events", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventListResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEvents) Patch(ctx context.Context, req *CalendarEventPatchReq, options ...lark.RequestOptionFunc) (*CalendarEventPatchResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPatch,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventPatchResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEvents) Search(ctx context.Context, req *CalendarEventSearchReq, options ...lark.RequestOptionFunc) (*CalendarEventSearchResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/search", []lark.AccessTokenType{lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventSearchResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEvents) Subscription(ctx context.Context, req *CalendarEventSubscriptionReq, options ...lark.RequestOptionFunc) (*CalendarEventSubscriptionResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/subscription", []lark.AccessTokenType{lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventSubscriptionResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEventAttendees) List(ctx context.Context, req *CalendarEventAttendeeListReq, options ...lark.RequestOptionFunc) (*CalendarEventAttendeeListResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodGet,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id/attendees", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventAttendeeListResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEventAttendees) BatchDelete(ctx context.Context, req *CalendarEventAttendeeBatchDeleteReq, options ...lark.RequestOptionFunc) (*CalendarEventAttendeeBatchDeleteResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id/attendees/batch_delete", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventAttendeeBatchDeleteResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEventAttendees) Create(ctx context.Context, req *CalendarEventAttendeeCreateReq, options ...lark.RequestOptionFunc) (*CalendarEventAttendeeCreateResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id/attendees", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventAttendeeCreateResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (c *calendarEventAttendeeChatMembers) List(ctx context.Context, req *CalendarEventAttendeeChatMemberListReq, options ...lark.RequestOptionFunc) (*CalendarEventAttendeeChatMemberListResp, error) {
	rawResp, err := c.app.SendRequestWithAccessTokenTypes(ctx, http.MethodGet,
		"/open-apis/calendar/v4/calendars/:calendar_id/events/:event_id/attendees/:attendee_id/chat_members", []lark.AccessTokenType{lark.AccessTokenTypeTenant, lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &CalendarEventAttendeeChatMemberListResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (f *freebusy) List(ctx context.Context, req *FreebusyListReq, options ...lark.RequestOptionFunc) (*FreebusyListResp, error) {
	rawResp, err := f.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/freebusy/list", []lark.AccessTokenType{lark.AccessTokenTypeTenant}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &FreebusyListResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (s *settings) GenerateCaldavConf(ctx context.Context, req *SettingGenerateCaldavConfReq, options ...lark.RequestOptionFunc) (*SettingGenerateCaldavConfResp, error) {
	rawResp, err := s.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/settings/generate_caldav_conf", []lark.AccessTokenType{lark.AccessTokenTypeUser}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &SettingGenerateCaldavConfResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (t *timeoffEvents) Delete(ctx context.Context, req *TimeoffEventDeleteReq, options ...lark.RequestOptionFunc) (*TimeoffEventDeleteResp, error) {
	rawResp, err := t.app.SendRequestWithAccessTokenTypes(ctx, http.MethodDelete,
		"/open-apis/calendar/v4/timeoff_events/:timeoff_event_id", []lark.AccessTokenType{lark.AccessTokenTypeTenant}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &TimeoffEventDeleteResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (t *timeoffEvents) Create(ctx context.Context, req *TimeoffEventCreateReq, options ...lark.RequestOptionFunc) (*TimeoffEventCreateResp, error) {
	rawResp, err := t.app.SendRequestWithAccessTokenTypes(ctx, http.MethodPost,
		"/open-apis/calendar/v4/timeoff_events", []lark.AccessTokenType{lark.AccessTokenTypeTenant}, req, options...)
	if err != nil {
		return nil, err
	}
	resp := &TimeoffEventCreateResp{RawResponse: rawResp}
	err = rawResp.JSONUnmarshalBody(resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}
