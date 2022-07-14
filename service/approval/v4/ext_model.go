package larkapproval

import larkevent "github.com/larksuite/oapi-sdk-go.v3/event"

type P1LeaveApprovalV4 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1LeaveApprovalV4Data `json:"event"`
}

func (m *P1LeaveApprovalV4) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1LeaveApprovalV4Data struct {
	AppID              string `json:"app_id,omitempty"`
	TenantKey          string `json:"tenant_key,omitempty"`
	Type               string `json:"type,omitempty"`
	InstanceCode       string `json:"instance_code,omitempty"`
	UserID             string `json:"user_id,omitempty"`
	OpenID             string `json:"open_id,omitempty"`
	OriginInstanceCode string `json:"origin_instance_code,omitempty"`
	StartTime          int64  `json:"start_time,omitempty"`
	EndTime            int64  `json:"end_time,omitempty"`

	LeaveFeedingArriveLate int64 `json:"leave_feeding_arrive_late,omitempty"`
	LeaveFeedingLeaveEarly int64 `json:"leave_feeding_leave_early,omitempty"`
	LeaveFeedingRestDaily  int64 `json:"leave_feeding_rest_daily,omitempty"`

	LeaveName      string                           `json:"leave_name,omitempty"`
	LeaveUnit      string                           `json:"leave_unit,omitempty"`
	LeaveStartTime string                           `json:"leave_start_time,omitempty"`
	LeaveEndTime   string                           `json:"leave_end_time,omitempty"`
	LeaveDetail    []string                         `json:"leave_detail,omitempty"`
	LeaveRange     []string                         `json:"leave_range,omitempty"`
	LeaveInterval  int64                            `json:"leave_interval,omitempty"`
	LeaveReason    string                           `json:"leave_reason,omitempty"`
	I18nResources  []*P1LeaveApprovalI18nResourceV4 `json:"i18n_resources,omitempty"`
}

type P1LeaveApprovalI18nResourceV4 struct {
	Locale    string            `json:"locale,omitempty"`     // 如: en_us
	IsDefault bool              `json:"is_default,omitempty"` // 如: true
	Texts     map[string]string `json:"texts,omitempty"`
}

type P1WorkApprovalV4 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1WorkApprovalV4Data `json:"event"`
}

func (m *P1WorkApprovalV4) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1WorkApprovalV4Data struct {
	AppID         string `json:"app_id,omitempty"`
	TenantKey     string `json:"tenant_key,omitempty"`
	Type          string `json:"type,omitempty"`
	InstanceCode  string `json:"instance_code,omitempty"`
	EmployeeID    string `json:"employee_id,omitempty"`
	OpenID        string `json:"open_id,omitempty"`
	StartTime     int64  `json:"start_time,omitempty"`
	EndTime       int64  `json:"end_time,omitempty"`
	WorkType      string `json:"work_type,omitempty"`
	WorkStartTime string `json:"work_start_time,omitempty"`
	WorkEndTime   string `json:"work_end_time,omitempty"`
	WorkInterval  int64  `json:"work_interval,omitempty"`
	WorkReason    string `json:"work_reason,omitempty"`
}

type P1ShiftApprovalV4 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1ShiftApprovalV4Data `json:"event"`
}

func (m *P1ShiftApprovalV4) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1ShiftApprovalV4Data struct {
	AppID        string `json:"app_id,omitempty"`
	TenantKey    string `json:"tenant_key,omitempty"`
	Type         string `json:"type,omitempty"`
	InstanceCode string `json:"instance_code,omitempty"`
	EmployeeID   string `json:"employee_id,omitempty"`
	OpenID       string `json:"open_id,omitempty"`
	StartTime    int64  `json:"start_time,omitempty"`
	EndTime      int64  `json:"end_time,omitempty"`
	ShiftTime    string `json:"shift_time,omitempty"`
	ReturnTime   string `json:"return_time,omitempty"`
	ShiftReason  string `json:"shift_reason,omitempty"`
}

type P1RemedyApprovalV4 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1RemedyApprovalV4Data `json:"event"`
}

func (m *P1RemedyApprovalV4) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1RemedyApprovalV4Data struct {
	AppID        string `json:"app_id,omitempty"`
	TenantKey    string `json:"tenant_key,omitempty"`
	Type         string `json:"type,omitempty"`
	InstanceCode string `json:"instance_code,omitempty"`
	EmployeeID   string `json:"employee_id,omitempty"`
	OpenID       string `json:"open_id,omitempty"`
	StartTime    int64  `json:"start_time,omitempty"`
	EndTime      int64  `json:"end_time,omitempty"`
	RemedyTime   string `json:"remedy_time,omitempty"`
	RemedyReason string `json:"remedy_reason,omitempty"`
}

type P1TripApprovalV4 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1TripApprovalV4Data `json:"event"`
}

func (m *P1TripApprovalV4) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1TripApprovalV4Data struct {
	AppID        string                      `json:"app_id,omitempty"`
	TenantKey    string                      `json:"tenant_key,omitempty"`
	Type         string                      `json:"type,omitempty"`
	InstanceCode string                      `json:"instance_code,omitempty"`
	EmployeeID   string                      `json:"employee_id,omitempty"`
	OpenID       string                      `json:"open_id,omitempty"`
	StartTime    int64                       `json:"start_time,omitempty"`
	EndTime      int64                       `json:"end_time,omitempty"`
	Schedules    []*P1TripApprovalScheduleV4 `json:"schedules,omitempty"`
	TripInterval int64                       `json:"trip_interval,omitempty"`
	TripReason   string                      `json:"trip_reason,omitempty"`
	TripPeers    []*P1TripApprovalTripPeerV4 `json:"trip_peers,omitempty"`
}

type P1TripApprovalScheduleV4 struct {
	TripStartTime  string `json:"trip_start_time,omitempty"`
	TripEndTime    string `json:"trip_end_time,omitempty"`
	TripInterval   int64  `json:"trip_interval,omitempty"`
	Departure      string `json:"departure,omitempty"`
	Destination    string `json:"destination,omitempty"`
	Transportation string `json:"transportation,omitempty"`
	TripType       string `json:"trip_type,omitempty"`
	Remark         string `json:"remark,omitempty"`
}

type P1TripApprovalTripPeerV4 struct {
	string `json:",omitempty"`
}

type P1OutApprovalV4 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1OutApprovalV4Data `json:"event"`
}

func (m *P1OutApprovalV4) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1OutApprovalV4Data struct {
	AppID         string                         `json:"app_id,omitempty"`
	I18nResources []*P1OutApprovalI18nResourceV4 `json:"i18n_resources,omitempty"`
	InstanceCode  string                         `json:"instance_code,omitempty"`
	OutImage      string                         `json:"out_image,omitempty"`
	OutInterval   int64                          `json:"out_interval,omitempty"`
	OutName       string                         `json:"out_name,omitempty"`
	OutReason     string                         `json:"out_reason,omitempty"`
	OutStartTime  string                         `json:"out_start_time,omitempty"`
	OutEndTime    string                         `json:"out_end_time,omitempty"`
	OutUnit       string                         `json:"out_unit,omitempty"`
	StartTime     int64                          `json:"start_time,omitempty"`
	EndTime       int64                          `json:"end_time,omitempty"`
	TenantKey     string                         `json:"tenant_key,omitempty"`
	Type          string                         `json:"type,omitempty"`
	OpenID        string                         `json:"open_id,omitempty"`
	UserID        string                         `json:"user_id,omitempty"`
}

type P1OutApprovalI18nResourceV4 struct {
	IsDefault bool              `json:"is_default,omitempty"`
	Locale    string            `json:"locale,omitempty"`
	Texts     map[string]string `json:"texts,omitempty"`
}
