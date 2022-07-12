package larkapplication

import larkevent "github.com/larksuite/oapi-sdk-go/event"

type P1AppOpenV6Data struct {
	AppID             string                        `json:"app_id,omitempty"`
	TenantKey         string                        `json:"tenant_key,omitempty"`
	Type              string                        `json:"type,omitempty"`
	Applicants        []*P1AppOpenApplicantV6       `json:"applicants,omitempty"`
	Installer         *P1AppOpenInstallerV6         `json:"installer,omitempty"`
	InstallerEmployee *P1AppOpenInstallerEmployeeV6 `json:"installer_employee,omitempty"`
}

type P1AppOpenV6 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1AppOpenV6Data `json:"event"`
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

func (m *P1AppOpenV6) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}
