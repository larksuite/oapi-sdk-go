package constants

const OAPIRootPath = "open-apis"
const (
	AppAccessTokenInternalUrlPath    string = "auth/v3/app_access_token/internal"
	AppAccessTokenIsvUrlPath         string = "auth/v3/app_access_token"
	TenantAccessTokenInternalUrlPath string = "auth/v3/tenant_access_token/internal"
	TenantAccessTokenIsvUrlPath      string = "auth/v3/tenant_access_token"
	ApplyAppTicketPath               string = "auth/v3/app_ticket/resend"
)

type IDType string

const (
	IDTypeOpen  IDType = "OpenID"
	IDTypeUnion IDType = "UnionID"
	IDTypeUser  IDType = "UserID"
)
