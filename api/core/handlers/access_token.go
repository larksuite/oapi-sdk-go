package handlers

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/token"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/store"
	"net/http"
	"time"
)

const expiryDelta = 3 * time.Minute

// get internal app access token
func getInternalAppAccessToken(ctx *core.Context) (*token.AppAccessToken, error) {
	accessToken := &token.AppAccessToken{}
	conf := core.GetConfigByCtx(ctx)
	req := request.NewRequestByAuth(AppAccessTokenInternalUrlPath, http.MethodPost,
		&token.InternalAccessTokenReq{
			AppID:     conf.GetAppSettings().AppID,
			AppSecret: conf.GetAppSettings().AppSecret,
		}, accessToken)
	err := send(ctx, req)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

// get internal tenant access token
func getInternalTenantAccessToken(ctx *core.Context) (*token.TenantAccessToken, error) {
	accessToken := &token.TenantAccessToken{}
	conf := core.GetConfigByCtx(ctx)
	req := request.NewRequestByAuth(TenantAccessTokenInternalUrlPath, http.MethodPost,
		&token.InternalAccessTokenReq{
			AppID:     conf.GetAppSettings().AppID,
			AppSecret: conf.GetAppSettings().AppSecret,
		}, accessToken)
	err := send(ctx, req)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

// get isv app access token
func getIsvAppAccessToken(ctx *core.Context) (*token.AppAccessToken, error) {
	appTicket, err := getAppTicket(ctx)
	if err != nil {
		return nil, err
	}
	if appTicket == "" {
		return nil, ErrAppTicketIsEmpty
	}
	conf := core.GetConfigByCtx(ctx)
	appAccessToken := &token.AppAccessToken{}
	req := request.NewRequestByAuth(AppAccessTokenIsvUrlPath, http.MethodPost,
		&token.ISVAppAccessTokenReq{
			AppID:     conf.GetAppSettings().AppID,
			AppSecret: conf.GetAppSettings().AppSecret,
			AppTicket: appTicket,
		}, appAccessToken)
	err = send(ctx, req)
	if err != nil {
		return nil, err
	}
	return appAccessToken, nil
}

func setAppAccessTokenToStore(ctx context.Context, appAccessToken *token.AppAccessToken) {
	conf := core.GetConfigByCtx(ctx)
	expire := time.Duration(appAccessToken.Expire)*time.Second - expiryDelta
	err := conf.GetStore().Put(ctx, store.AppAccessTokenKey(conf.GetAppSettings().AppID), appAccessToken.AppAccessToken, expire)
	if err != nil {
		conf.GetLogger().Warn(ctx, err)
	}
}

// get isv tenant access token
func getIsvTenantAccessToken(ctx *core.Context) (*token.AppAccessToken, *token.TenantAccessToken, error) {
	appAccessToken, err := getIsvAppAccessToken(ctx)
	if err != nil {
		return nil, nil, err
	}
	info := request.GetInfoByCtx(ctx)
	tenantAccessToken := &token.TenantAccessToken{}
	req := request.NewRequestByAuth(TenantAccessTokenIsvUrlPath, http.MethodPost,
		&token.ISVTenantAccessTokenReq{
			AppAccessToken: appAccessToken.AppAccessToken,
			TenantKey:      info.TenantKey,
		}, tenantAccessToken)
	err = send(ctx, req)
	if err != nil {
		return appAccessToken, nil, err
	}
	return appAccessToken, tenantAccessToken, nil
}

func setTenantAccessTokenToStore(ctx context.Context, tenantAccessToken *token.TenantAccessToken) {
	info := request.GetInfoByCtx(ctx)
	conf := core.GetConfigByCtx(ctx)
	expire := time.Duration(tenantAccessToken.Expire)*time.Second - expiryDelta
	err := conf.GetStore().Put(ctx, store.TenantAccessTokenKey(conf.GetAppSettings().AppID, info.TenantKey), tenantAccessToken.TenantAccessToken, expire)
	if err != nil {
		conf.GetLogger().Warn(ctx, err)
	}
}

func setAuthorizationToHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func cloneHTTPRequest(req *http.Request) *http.Request {
	convertedRequest := new(http.Request)
	*convertedRequest = *req
	convertedRequest.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		convertedRequest.Header[k] = append([]string(nil), s...)
	}
	return convertedRequest
}

func send(ctx *core.Context, req *request.Request) error {
	Handle(ctx, req)
	return req.Err
}

func getAppTicket(ctx *core.Context) (string, error) {
	conf := core.GetConfigByCtx(ctx)
	return conf.GetStore().Get(ctx, store.AppTicketKey(conf.GetAppSettings().AppID))
}

func setAppAccessToken(ctx *core.Context, req *http.Request) (*http.Request, error) {
	convertedRequest := cloneHTTPRequest(req)
	info := request.GetInfoByCtx(ctx)
	conf := core.GetConfigByCtx(ctx)
	// from store get app access token
	if !info.Retryable {
		tok, err := conf.GetStore().Get(ctx, store.AppAccessTokenKey(conf.GetAppSettings().AppID))
		if err != nil {
			return nil, err
		}
		if tok != "" {
			setAuthorizationToHeader(convertedRequest, tok)
			return convertedRequest, nil
		}
	}
	// from api get app access token
	var appAccessToken *token.AppAccessToken
	var err error
	if conf.GetAppSettings().AppType == constants.AppTypeInternal {
		appAccessToken, err = getInternalAppAccessToken(ctx)
	} else {
		appAccessToken, err = getIsvAppAccessToken(ctx)
	}
	if err != nil {
		return nil, err
	}
	setAppAccessTokenToStore(ctx, appAccessToken)
	setAuthorizationToHeader(convertedRequest, appAccessToken.AppAccessToken)
	return convertedRequest, nil
}

func setTenantAccessToken(ctx *core.Context, req *http.Request) (*http.Request, error) {
	convertedRequest := cloneHTTPRequest(req)
	info := request.GetInfoByCtx(ctx)
	conf := core.GetConfigByCtx(ctx)
	// from store get tenant access token
	if !info.Retryable {
		tok, err := conf.GetStore().Get(ctx, store.TenantAccessTokenKey(conf.GetAppSettings().AppID, info.TenantKey))
		if err != nil {
			return nil, err
		}
		if tok != "" {
			setAuthorizationToHeader(convertedRequest, tok)
			return convertedRequest, nil
		}
	}
	// from api get tenant access token
	var tenantAccessToken *token.TenantAccessToken
	var appAccessToken *token.AppAccessToken
	var err error
	if conf.GetAppSettings().AppType == constants.AppTypeInternal {
		tenantAccessToken, err = getInternalTenantAccessToken(ctx)
	} else {
		appAccessToken, tenantAccessToken, err = getIsvTenantAccessToken(ctx)
		if appAccessToken != nil {
			setAppAccessTokenToStore(ctx, appAccessToken)
		}
	}
	if err != nil {
		return nil, err
	}
	setTenantAccessTokenToStore(ctx, tenantAccessToken)
	setAuthorizationToHeader(convertedRequest, tenantAccessToken.TenantAccessToken)
	return convertedRequest, nil
}

func setUserAccessToken(ctx *core.Context, req *http.Request) (*http.Request, error) {
	convertedRequest := cloneHTTPRequest(req)
	info := request.GetInfoByCtx(ctx)
	setAuthorizationToHeader(convertedRequest, info.UserAccessToken)
	return convertedRequest, nil
}
