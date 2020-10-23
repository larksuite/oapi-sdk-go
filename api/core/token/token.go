package token

import (
	"time"
)

type GetISVTenantAccessTokenReq struct {
	AppAccessToken string `json:"app_access_token"`
	TenantKey      string `json:"tenant_key"`
}

type TenantAccessToken struct {
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

type GetInternalAccessTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type GetISVAppAccessTokenReq struct {
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

type UserToken struct {
	Id                string
	AppID             string
	AccessToken       string
	ExpiresIn         int
	ExpireTime        time.Time
	OpenID            string
	UnionID           string
	UserID            string
	TenantKey         string
	RefreshExpiresIn  int
	RefreshExpireTime time.Time
	RefreshToken      string
	TokenType         string
	Error             string
}
