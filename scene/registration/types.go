package registration

const (
	StatusPolling        = "polling"
	StatusSlowDown       = "slow_down"
	StatusDomainSwitched = "domain_switched"
)

type QRCodeInfo struct {
	URL      string
	ExpireIn int
}

type StatusChangeInfo struct {
	Status   string
	Interval int
}

type UserInfo struct {
	OpenID      string
	TenantBrand string
}

type RegisterAppResult struct {
	ClientID     string
	ClientSecret string
	UserInfo     *UserInfo
}

type Options struct {
	Source         string
	Domain         string
	LarkDomain     string
	OnQRCode       func(info *QRCodeInfo)
	OnStatusChange func(info *StatusChangeInfo)
}

type beginResponse struct {
	DeviceCode              string `json:"device_code"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	VerificationURI         string `json:"verification_uri"`
	UserCode                string `json:"user_code"`
	Interval                int    `json:"interval"`
	ExpireIn                int    `json:"expire_in"`
}

type pollResponse struct {
	ClientID     string    `json:"client_id,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	UserInfo     *userInfo `json:"user_info,omitempty"`
	Error        string    `json:"error,omitempty"`
	ErrorDesc    string    `json:"error_description,omitempty"`
}

type userInfo struct {
	OpenID      string `json:"open_id,omitempty"`
	TenantBrand string `json:"tenant_brand,omitempty"`
}
