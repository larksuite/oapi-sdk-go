package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/larksuite/oapi-sdk-go/httpclient"
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

func validate(config *Config, option *RequestOption, accessTokenType AccessTokenType) error {
	if config.AppId == "" {
		return &IllegalParamError{msg: "AppId is empty"}
	}

	if config.AppSecret == "" {
		return &IllegalParamError{msg: "AppSecret is empty"}
	}

	if !config.EnableTokenCache && option.UserAccessToken == "" && option.TenantAccessToken == "" && option.AppAccessToken == "" {
		return &IllegalParamError{msg: "accessToken is empty"}
	}

	if config.AppType == AppTypeMarketplace && accessTokenType == AccessTokenTypeTenant && option.TenantKey == "" {
		return &IllegalParamError{msg: "tenant key is empty"}
	}

	if accessTokenType == AccessTokenTypeUser && option.UserAccessToken == "" {
		return &IllegalParamError{msg: "user access token is empty"}
	}

	if option.Header != nil {
		if option.Header.Get(httpHeaderKeyRequestId) != "" {
			return &IllegalParamError{msg: fmt.Sprintf("use %s as header key is not allowed", httpHeaderKeyRequestId)}
		}
		if option.Header.Get(httpHeaderRequestId) != "" {
			return &IllegalParamError{msg: fmt.Sprintf("use %s as header key is not allowed", httpHeaderRequestId)}
		}
		if option.Header.Get(httpHeaderKeyLogId) != "" {
			return &IllegalParamError{msg: fmt.Sprintf("use %s as header key is not allowed", httpHeaderKeyLogId)}
		}
	}

	return nil
}

func doSend(ctx context.Context, rawRequest *http.Request, httpClient httpclient.HttpClient) (*RawResponse, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(rawRequest)
	if err != nil {
		if er, ok := err.(*url.Error); ok {
			if er.Timeout() {
				return nil, &ClientTimeoutError{msg: er.Error()}
			}

			if e, ok := er.Err.(*net.OpError); ok && e.Op == "dial" {
				return nil, &DialFailedError{msg: er.Error()}
			}
		}
		return nil, err
	}

	if resp.StatusCode == http.StatusGatewayTimeout {
		return nil, &ServerTimeoutError{msg: "server time out error "}
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

type applyAppTicketReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func SendRequest(ctx context.Context, config *Config, httpMethod string, httpPath string,
	accessTokenTypes []AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	option := &RequestOption{}
	for _, optionFunc := range options {
		optionFunc(option)
	}

	accessTokenType := determineTokenType(accessTokenTypes, option, config.EnableTokenCache)
	err := validate(config, option, accessTokenType)
	if err != nil {
		return nil, err
	}

	rawResp, err := doSendRequest(ctx, config, httpMethod, httpPath, accessTokenType, input, option)
	return rawResp, err
}

func doSendRequest(ctx context.Context, config *Config, httpMethod string, httpPath string,
	accessTokenType AccessTokenType, input interface{}, option *RequestOption) (*RawResponse, error) {
	var rawResp *RawResponse
	var errResult error
	for i := 0; i < 2; i++ {
		req, err := reqTranslator.translate(ctx, input, accessTokenType, config, httpMethod, httpPath, option)
		if err != nil {
			return nil, err
		}

		if config.LogReqRespInfoAtDebugLevel {
			config.Logger.Debug(ctx, fmt.Sprintf("req:%v", req))
		} else {
			config.Logger.Debug(ctx, fmt.Sprintf("req:%s,%s", httpMethod, httpPath))
		}
		rawResp, err = doSend(ctx, req, config.HttpClient)
		if config.LogReqRespInfoAtDebugLevel {
			config.Logger.Debug(ctx, fmt.Sprintf("resp:%v", rawResp))
		}
		_, isDialError := err.(*DialFailedError)
		if err != nil && !isDialError {
			return nil, err
		}
		errResult = err
		if isDialError {
			continue
		}

		fileDownloadSuccess := option.FileDownload && rawResp.StatusCode == http.StatusOK
		if fileDownloadSuccess || !strings.Contains(rawResp.Header.Get(contentTypeHeader), contentTypeJson) {
			break
		}

		codeError := &CodeError{}
		err = json.Unmarshal(rawResp.RawBody, codeError)
		if err != nil {
			return nil, err
		}

		code := codeError.Code
		if code == errCodeAppTicketInvalid {
			applyAppTicket(ctx, config)
		}

		if accessTokenType == accessTokenTypeNone {
			break
		}

		if !config.EnableTokenCache {
			break
		}

		if code != errCodeAccessTokenInvalid && code != errCodeAppAccessTokenInvalid &&
			code != errCodeTenantAccessTokenInvalid {
			break
		}
	}

	if errResult != nil {
		return nil, errResult
	}
	return rawResp, nil
}

func GetConfig(appId, appSecret string) *Config {
	config := &Config{
		Domain:           "https://open.feishu.cn",
		AppType:          AppTypeSelfBuilt,
		EnableTokenCache: true,
		AppSecret:        appSecret,
		AppId:            appId,
	}

	NewLogger(config)
	NewCache(config)

	return config
}

//
//func Post(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodPost, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func Get(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodGet, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func Head(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodHead, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func Put(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodPut, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func Patch(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodPatch, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func Delete(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodDelete, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func Connect(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodConnect, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func SendOptions(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodOptions, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
//
//func SendTrace(ctx context.Context, config *Config, httpPath string,
//	body interface{}, accessTokeType AccessTokenType, options ...RequestOptionFunc) (*RawResponse, error) {
//	return SendRequest(ctx, config, http.MethodTrace, httpPath, []AccessTokenType{accessTokeType}, body, options...)
//}
