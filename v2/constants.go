package lark

import "errors"

const contentTypeHeader = "Content-Type"
const contentTypeJson = "application/json"
const defaultContentType = contentTypeJson + "; charset=utf-8"

const (
	httpHeaderKeyRequestID = "X-Request-Id"
	httpHeaderKeyLogID     = "X-Tt-Logid"
)

type AppType string

const (
	AppTypeCustom      AppType = "Custom App"
	AppTypeMarketplace AppType = "Marketplace App"
)

type Domain string

const (
	DomainFeiShu    Domain = "https://open.feishu.cn"
	DomainLarkSuite Domain = "https://open.larksuite.com"
)

type CallbackType string

const (
	callbackTypeEvent     CallbackType = "event_callback"
	callbackTypeChallenge CallbackType = "url_verification"
)

type AccessTokenType string

const (
	accessTokenTypeNone   AccessTokenType = "none_access_token"
	AccessTokenTypeApp    AccessTokenType = "app_access_token"
	AccessTokenTypeTenant AccessTokenType = "tenant_access_token"
	AccessTokenTypeUser   AccessTokenType = "user_access_token"
)

var (
	ErrTenantKeyIsEmpty          = errors.New("tenant key is empty")
	ErrUserAccessTokenKeyIsEmpty = errors.New("user access token is empty")
	ErrAppTicketIsEmpty          = errors.New("app ticket is empty")
)

const (
	appAccessTokenInternalUrlPath    string = "/open-apis/auth/v3/app_access_token/internal"
	appAccessTokenUrlPath            string = "/open-apis/auth/v3/app_access_token"
	tenantAccessTokenInternalUrlPath string = "/open-apis/auth/v3/tenant_access_token/internal"
	tenantAccessTokenUrlPath         string = "/open-apis/auth/v3/tenant_access_token"
	applyAppTicketPath               string = "/open-apis/auth/v3/app_ticket/resend"
)

const (
	errCodeOk                       = 0
	errCodeAppTicketInvalid         = 10012
	errCodeAccessTokenInvalid       = 99991671
	errCodeAppAccessTokenInvalid    = 99991664
	errCodeTenantAccessTokenInvalid = 99991663
	ErrCodeUserAccessTokenInvalid   = 99991668
	ErrCodeUserRefreshTokenInvalid  = 99991669
)

const version = "1.0.0"
