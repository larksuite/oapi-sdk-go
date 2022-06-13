package core

import (
	"context"
	"errors"
	"net/http"
)

var reqTranslator ReqTranslator

func determineTokenType(accessTokenTypes []AccessTokenType, option *RequestOption, enableTokenCache bool) AccessTokenType {
	if !enableTokenCache {
		if option.UserAccessToken != "" {
			return AccessTokenTypeUser
		}
		if option.TenantAccessToken != "" {
			return AccessTokenTypeTenant
		}
		if option.AppAccessToken != "" {
			return AccessTokenTypeApp
		}

		return accessTokenTypeNone
	}
	accessibleTokenTypeSet := make(map[AccessTokenType]struct{})
	accessTokenType := accessTokenTypes[0]
	for _, t := range accessTokenTypes {
		if t == AccessTokenTypeTenant {
			accessTokenType = t // default
		}
		accessibleTokenTypeSet[t] = struct{}{}
	}
	if option.TenantKey != "" {
		if _, ok := accessibleTokenTypeSet[AccessTokenTypeTenant]; ok {
			accessTokenType = AccessTokenTypeTenant
		}
	}
	if option.UserAccessToken != "" {
		if _, ok := accessibleTokenTypeSet[AccessTokenTypeUser]; ok {
			accessTokenType = AccessTokenTypeUser
		}
	}

	return accessTokenType
}

func valid(config *Config, option *RequestOption, accessTokenType AccessTokenType) error {
	if config.EnableTokenCache == false && option.UserAccessToken == "" && option.TenantAccessToken == "" && option.AppAccessToken == "" {
		return errors.New("accessToken is empty")
	}

	if config.AppType == AppTypeMarketplace && accessTokenType == AccessTokenTypeTenant && option.TenantKey == "" {
		return errors.New("tenant key is empty")
	}

	if accessTokenType == AccessTokenTypeUser && option.UserAccessToken == "" {
		return errors.New("user access token is empty")
	}

	return nil
}

func doSend(rawRequest *http.Request, httpClient *http.Client) (*RawResponse, error) {
	resp, err := httpClient.Do(rawRequest)
	if err != nil {
		return nil, err
	}
	body, err := readResponse(resp)
	if err != nil {
		return nil, err
	}
	return &RawResponse{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		RawBody:    body,
	}, nil
}
func SendRequest(ctx context.Context, config *Config, httpMethod string, httpPath string,
	accessTokenTypes []AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {

	option := &RequestOption{}
	for _, optionFunc := range options {
		optionFunc(option)
	}

	accessTokenType := determineTokenType(accessTokenTypes, option, config.EnableTokenCache)

	err := valid(config, option, accessTokenType)
	if err != nil {
		return nil, err
	}

	req, err := reqTranslator.translate(ctx, input, accessTokenType, config, httpMethod, httpPath, option)
	if err != nil {
		return nil, err
	}

	rawResp, err := doSend(req, config.HttpClient)

	return rawResp, err
}

func SendPost(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodPost, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendGet(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodGet, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendHead(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodHead, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendPut(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodPut, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendPatch(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodPatch, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendDelete(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodDelete, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendConnect(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodConnect, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendOptions(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodOptions, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}

func SendTrace(ctx context.Context, config *Config, httpPath string,
	accessTokeType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return SendRequest(ctx, config, http.MethodTrace, httpPath, []AccessTokenType{accessTokeType}, input, options...)
}
