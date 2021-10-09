package lark

const userAgentHeader = "User-Agent"
const contentTypeHeader = "Content-Type"
const contentTypeJson = "application/json"
const defaultContentType = contentTypeJson + "; charset=utf-8"

const (
	httpHeaderKeyRequestId = "X-Request-Id"
	httpHeaderKeyLogId     = "X-Tt-Logid"
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

type webhookType string

const (
	webhookTypeChallenge webhookType = "url_verification"
)

type AccessTokenType string

const (
	accessTokenTypeNone   AccessTokenType = "none_access_token"
	AccessTokenTypeApp    AccessTokenType = "app_access_token"
	AccessTokenTypeTenant AccessTokenType = "tenant_access_token"
	AccessTokenTypeUser   AccessTokenType = "user_access_token"
)

const (
	appAccessTokenInternalUrlPath    string = "/open-apis/auth/v3/app_access_token/internal"
	appAccessTokenUrlPath            string = "/open-apis/auth/v3/app_access_token"
	tenantAccessTokenInternalUrlPath string = "/open-apis/auth/v3/tenant_access_token/internal"
	tenantAccessTokenUrlPath         string = "/open-apis/auth/v3/tenant_access_token"
	applyAppTicketPath               string = "/open-apis/auth/v3/app_ticket/resend"
)

const (
	errCodeAppTicketInvalid         = 10012
	errCodeAccessTokenInvalid       = 99991671
	errCodeAppAccessTokenInvalid    = 99991664
	errCodeTenantAccessTokenInvalid = 99991663
)

const (
	larkRequestNonce     = "X-Lark-Request-Nonce"
	larkRequestTimestamp = "X-Lark-Request-Timestamp"
	larkSignature        = "X-Lark-Signature"
)

const webhookResponseFormat = `{"msg":"%s"}`
const challengeResponseFormat = `{"challenge":"%s"}`
