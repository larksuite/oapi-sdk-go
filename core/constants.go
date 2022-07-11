package larkcore

import (
	"time"
)

const defaultContentType = contentTypeJson + "; charset=utf-8"
const userAgentHeader = "User-Agent"

const (
	HttpHeaderKeyRequestId = "X-Request-Id"
	httpHeaderRequestId    = "Request-Id"
	HttpHeaderKeyLogId     = "X-Tt-Logid"
	contentTypeHeader      = "Content-Type"
	contentTypeJson        = "application/json"
	customRequestId        = "Oapi-Sdk-Request-Id"
)

type AppType string

const (
	AppTypeSelfBuilt   AppType = "SelfBuilt"
	AppTypeMarketplace AppType = "Marketplace"
)

const (
	appAccessTokenInternalUrlPath    string = "/open-apis/auth/v3/app_access_token/internal"
	appAccessTokenUrlPath            string = "/open-apis/auth/v3/app_access_token"
	tenantAccessTokenInternalUrlPath string = "/open-apis/auth/v3/tenant_access_token/internal"
	tenantAccessTokenUrlPath         string = "/open-apis/auth/v3/tenant_access_token"
	applyAppTicketPath               string = "/open-apis/auth/v3/app_ticket/resend"
)

type AccessTokenType string

const (
	accessTokenTypeNone   AccessTokenType = "none_access_token"
	AccessTokenTypeApp    AccessTokenType = "app_access_token"
	AccessTokenTypeTenant AccessTokenType = "tenant_access_token"
	AccessTokenTypeUser   AccessTokenType = "user_access_token"
)

const (
	appTicketKeyPrefix         = "app_ticket"
	appAccessTokenKeyPrefix    = "app_access_token"
	tenantAccessTokenKeyPrefix = "tenant_access_token"
)
const expiryDelta = 3 * time.Minute
const (
	errCodeAppTicketInvalid         = 10012
	errCodeAccessTokenInvalid       = 99991671
	errCodeAppAccessTokenInvalid    = 99991664
	errCodeTenantAccessTokenInvalid = 99991663
)
