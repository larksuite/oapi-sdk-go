package constants

const ContentType = "Content-Type"
const ContentTypeJson = "application/json"
const DefaultContentType = ContentTypeJson + "; charset=utf-8"

const (
	HTTPHeader             = "HTTP-Header"
	HTTPHeaderKeyRequestID = "X-Request-Id"
	HTTPHeaderKeyLogID     = "X-Tt-Logid"
	HTTPKeyStatusCode      = "http_status_code"
)

type AppType string

const (
	AppTypeISV      AppType = "isv"      // isv app
	AppTypeInternal AppType = "internal" // internal app
)

type Domain string

const (
	DomainFeiShu    Domain = "https://open.feishu.cn"
	DomainLarkSuite Domain = "https://open.larksuite.com"
)

type CallbackType string

const (
	CallbackTypeEvent     CallbackType = "event_callback"
	CallbackTypeChallenge CallbackType = "url_verification"
)
