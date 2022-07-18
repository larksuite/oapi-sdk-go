package larkcore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var tokenManager TokenManager = TokenManager{cache: cache}

type TokenManager struct {
	cache Cache
}

func (m *TokenManager) getAppAccessToken(ctx context.Context, config *Config) (string, error) {
	token, err := m.get(ctx, appAccessTokenKey(config.AppId))
	if err != nil {
		return "", err
	}

	appType := config.AppType
	if token == "" {
		if appType == AppTypeSelfBuilt {
			token, err = m.getCustomAppAccessTokenThenCache(ctx, config)
			if err != nil {
				return "", err
			}
			return token, nil
		} else {
			token, err = m.getMarketplaceAppAccessTokenThenCache(ctx, config)
			if err != nil {
				return "", err
			}
			return token, nil
		}
	}
	return token, nil
}

func (m *TokenManager) getTenantAccessToken(ctx context.Context, config *Config, tenantKey string) (string, error) {
	token, err := m.get(ctx, tenantAccessTokenKey(config.AppId, tenantKey))
	if err != nil {
		return "", err
	}

	if token == "" {
		if config.AppType == AppTypeSelfBuilt {
			token, err = m.getCustomTenantAccessTokenThenCache(ctx, config, tenantKey)
			if err != nil {
				return "", err
			}
			return token, nil
		} else {
			token, err = m.getMarketplaceTenantAccessTokenThenCache(ctx, config, tenantKey)
			if err != nil {
				return "", err
			}
			return token, nil
		}
	}
	return token, nil
}

func (m *TokenManager) set(ctx context.Context, key, value string, ttl time.Duration) error {
	return m.cache.Set(ctx, key, value, ttl)
}

func (m *TokenManager) get(ctx context.Context, tokenKey string) (string, error) {
	token, err := m.cache.Get(ctx, tokenKey)
	return token, err
}

type internalAccessTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}
type appAccessTokenResp struct {
	CodeError
	Expire         int    `json:"expire"`
	AppAccessToken string `json:"app_access_token"`
}
type tenantAccessTokenResp struct {
	CodeError
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}
type marketplaceAppAccessTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	AppTicket string `json:"app_ticket"`
}

type marketplaceTenantAccessTokenReq struct {
	AppAccessToken string `json:"app_access_token"`
	TenantKey      string `json:"tenant_key"`
}

func appAccessTokenKey(appID string) string {
	return fmt.Sprintf("%s-%s", appAccessTokenKeyPrefix, appID)
}

func tenantAccessTokenKey(appID, tenantKey string) string {
	return fmt.Sprintf("%s-%s-%s", tenantAccessTokenKeyPrefix, appID, tenantKey)
}
func (m *TokenManager) getCustomAppAccessTokenThenCache(ctx context.Context, config *Config) (string, error) {
	rawResp, err := Request(ctx, &ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    appAccessTokenInternalUrlPath,
		Body: &internalAccessTokenReq{
			AppID:     config.AppId,
			AppSecret: config.AppSecret,
		},
		SupportedAccessTokenTypes: []AccessTokenType{accessTokenTypeNone},
	}, config)

	if err != nil {
		return "", err
	}
	resp := &appAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		config.Logger.Warn(ctx, fmt.Sprintf("custom app appAccessToken cache, err:%v", Prettify(resp)))
		return "", resp.CodeError
	}
	expire := time.Duration(resp.Expire)*time.Second - expiryDelta
	err = m.set(ctx, appAccessTokenKey(config.AppId), resp.AppAccessToken, expire)
	if err != nil {
		config.Logger.Warn(ctx, fmt.Sprintf("custom app appAccessToken save cache, err:%v", err))
	}
	return resp.AppAccessToken, err
}

func (m *TokenManager) getCustomTenantAccessTokenThenCache(ctx context.Context, config *Config, tenantKey string) (string, error) {
	rawResp, err := Request(ctx, &ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    tenantAccessTokenInternalUrlPath,
		Body: &internalAccessTokenReq{
			AppID:     config.AppId,
			AppSecret: config.AppSecret,
		},
		SupportedAccessTokenTypes: []AccessTokenType{accessTokenTypeNone},
	}, config)

	if err != nil {
		return "", err
	}
	tenantAccessTokenResp := &tenantAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, tenantAccessTokenResp)
	if err != nil {
		return "", err
	}
	if tenantAccessTokenResp.Code != 0 {
		config.Logger.Warn(ctx, fmt.Sprintf("custom app tenantAccessToken cache, err:%v", Prettify(tenantAccessTokenResp)))
		return "", tenantAccessTokenResp.CodeError
	}
	expire := time.Duration(tenantAccessTokenResp.Expire)*time.Second - expiryDelta
	err = m.cache.Set(ctx, tenantAccessTokenKey(config.AppId, tenantKey), tenantAccessTokenResp.TenantAccessToken, expire)
	if err != nil {
		config.Logger.Warn(ctx, fmt.Sprintf("custom app tenantAccessToken save cache, err:%v", err))
	}
	return tenantAccessTokenResp.TenantAccessToken, err
}

var ErrAppTicketIsEmpty = errors.New("app ticket is empty")

func (m *TokenManager) getMarketplaceAppAccessTokenThenCache(ctx context.Context, config *Config) (string, error) {
	appTicket, err := appTicketManager.Get(ctx, config)
	if err != nil {
		return "", err
	}
	if appTicket == "" {
		return "", ErrAppTicketIsEmpty
	}
	rawResp, err := Request(ctx, &ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    appAccessTokenUrlPath,
		Body: &marketplaceAppAccessTokenReq{
			AppID:     config.AppId,
			AppSecret: config.AppSecret,
			AppTicket: appTicket,
		},
		SupportedAccessTokenTypes: []AccessTokenType{accessTokenTypeNone},
	}, config)

	if err != nil {
		return "", err
	}
	appAccessTokenResp := &appAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, appAccessTokenResp)
	if err != nil {
		config.Logger.Warn(ctx, fmt.Sprintf("marketplace app appAccessToken cache, err:%v", Prettify(appAccessTokenResp)))
		return "", err
	}
	if appAccessTokenResp.Code != 0 {
		return "", appAccessTokenResp.CodeError
	}
	expire := time.Duration(appAccessTokenResp.Expire)*time.Second - expiryDelta
	err = m.set(ctx, appAccessTokenKey(config.AppId), appAccessTokenResp.AppAccessToken, expire)
	if err != nil {
		config.Logger.Warn(ctx, fmt.Sprintf("marketplace app appAccessToken save cache, err:%v", err))
	}
	return appAccessTokenResp.AppAccessToken, err
}

// get marketplace tenant access token
func (m *TokenManager) getMarketplaceTenantAccessTokenThenCache(ctx context.Context, config *Config, tenantKey string) (string, error) {
	appAccessToken, err := m.getMarketplaceAppAccessTokenThenCache(ctx, config)
	if err != nil {
		return "", err
	}
	rawResp, err := Request(ctx, &ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    tenantAccessTokenUrlPath,
		Body: &marketplaceTenantAccessTokenReq{
			AppAccessToken: appAccessToken,
			TenantKey:      tenantKey,
		},
		SupportedAccessTokenTypes: []AccessTokenType{accessTokenTypeNone},
	}, config)

	if err != nil {
		return "", err
	}
	tenantAccessTokenResp := &tenantAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, tenantAccessTokenResp)
	if err != nil {
		config.Logger.Warn(ctx, fmt.Sprintf("marketplace app tenantAccessToken cache, err:%v", Prettify(tenantAccessTokenResp)))
		return "", err
	}
	if tenantAccessTokenResp.Code != 0 {
		return "", tenantAccessTokenResp.CodeError
	}
	expire := time.Duration(tenantAccessTokenResp.Expire)*time.Second - expiryDelta
	err = m.set(ctx, tenantAccessTokenKey(config.AppId, tenantKey), tenantAccessTokenResp.TenantAccessToken, expire)
	if err != nil {
		config.Logger.Warn(ctx, fmt.Sprintf("market app tenantAccessToken save cache, err:%v", err))
	}
	return tenantAccessTokenResp.TenantAccessToken, err
}
