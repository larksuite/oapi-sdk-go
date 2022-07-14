package larkapplication

import larkevent "github.com/larksuite/oapi-sdk-go.v3/event"

type P1OrderPaidV6 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1OrderPaidV6Data `json:"event"`
}

func (m *P1OrderPaidV6) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1OrderPaidV6Data struct {
	Type          string `json:"type,omitempty"`
	AppID         string `json:"app_id,omitempty"`
	OrderID       string `json:"order_id,omitempty"`
	PricePlanID   string `json:"price_plan_id,omitempty"`
	PricePlanType string `json:"price_plan_type,omitempty"`
	BuyCount      int64  `json:"buy_count,omitempty"`
	Seats         int64  `json:"seats,omitempty"`
	CreateTime    string `json:"create_time,omitempty"`
	PayTime       string `json:"pay_time,omitempty"`
	BuyType       string `json:"buy_type,omitempty"`
	SrcOrderID    string `json:"src_order_id,omitempty"`
	OrderPayPrice int64  `json:"order_pay_price,omitempty"`
	TenantKey     string `json:"tenant_key,omitempty"`
}

type P1AppUninstalledV6 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1AppUninstalledV6Data `json:"event"`
}

type P1AppUninstalledV6Data struct {
	AppID     string `json:"app_id,omitempty"`
	TenantKey string `json:"tenant_key,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (m *P1AppUninstalledV6) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1AppOpenV6 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1AppOpenV6Data `json:"event"`
}

func (m *P1AppOpenV6) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1AppOpenApplicantV6 struct {
	OpenID string `json:"open_id,omitempty"`
}

type P1AppOpenInstallerV6 struct {
	OpenID string `json:"open_id,omitempty"`
}

type P1AppOpenInstallerEmployeeV6 struct {
	OpenID string `json:"open_id,omitempty"`
}

type P1AppOpenV6Data struct {
	AppID             string                        `json:"app_id,omitempty"`
	TenantKey         string                        `json:"tenant_key,omitempty"`
	Type              string                        `json:"type,omitempty"`
	Applicants        []*P1AppOpenApplicantV6       `json:"applicants,omitempty"`
	Installer         *P1AppOpenInstallerV6         `json:"installer,omitempty"`
	InstallerEmployee *P1AppOpenInstallerEmployeeV6 `json:"installer_employee,omitempty"`
}

type P1AppStatusChangedV6 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1AppStatusChangedV6Data `json:"event"`
}

func (m *P1AppStatusChangedV6) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1AppStatusChangedV6Data struct {
	AppID     string                       `json:"app_id,omitempty"`
	TenantKey string                       `json:"tenant_key,omitempty"`
	Type      string                       `json:"type,omitempty"`
	Status    string                       `json:"status,omitempty"`
	Operator  *P1AppStatusChangeOperatorV6 `json:"operator,omitempty"`
}

type P1AppStatusChangeOperatorV6 struct {
	OpenID  string `json:"open_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	UnionId string `json:"union_id,omitempty"`
}
