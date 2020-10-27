package token

type ISVTenantAccessTokenReq struct {
	AppAccessToken string `json:"app_access_token"`
	TenantKey      string `json:"tenant_key"`
}

type TenantAccessToken struct {
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

type InternalAccessTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type ISVAppAccessTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	AppTicket string `json:"app_ticket"`
}

type AppAccessToken struct {
	Expire         int    `json:"expire"`
	AppAccessToken string `json:"app_access_token"`
}

type ApplyAppTicketReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}
