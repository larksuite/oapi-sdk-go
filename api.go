package lark

import (
	"github.com/larksuite/oapi-sdk-go/api"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
)

type NoData = response.NoData

type AccessTokenType = request.AccessTokenType

const (
	AccessTokenTypeNone   = request.AccessTokenTypeNone
	AccessTokenTypeApp    = request.AccessTokenTypeApp
	AccessTokenTypeTenant = request.AccessTokenTypeTenant
	AccessTokenTypeUser   = request.AccessTokenTypeUser
)

type APIRequest = request.Request

type APIError = response.Error

type APIRequestOpt = request.OptFn

var (
	SetUserAccessToken = request.SetUserAccessToken
	SetTenantKey       = request.SetTenantKey
	SetQueryParams     = request.SetQueryParams
	SetPathParams      = request.SetPathParams
	SetNotDataField    = request.SetNotDataField
	SetResponseStream  = request.SetResponseStream
	NeedHelpDeskAuth   = request.NeedHelpDeskAuth
)

type FormData = request.FormData
type FormDataFile = request.File

var (
	NewFormDataFile = request.NewFile
	NewFormData     = request.NewFormData
)

func NewAPIRequest(httpPath, httpMethod string, accessTokenType AccessTokenType,
	input interface{}, output interface{}, opts ...APIRequestOpt) *APIRequest {
	req := request.NewRequestWithNative(httpPath, httpMethod, accessTokenType, input, output, opts...)
	return req
}

func NewAPIRequestWithMultiToken(httpPath, httpMethod string, accessTokenTypes []AccessTokenType,
	input interface{}, output interface{}, opts ...APIRequestOpt) *APIRequest {
	req := request.NewRequest(httpPath, httpMethod, accessTokenTypes, input, output, opts...)
	return req
}

func SendAPIRequest(ctx *core.Context, conf Config, req *APIRequest) error {
	return api.Send(ctx, conf, req)
}
