package request

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/constants"
	"github.com/larksuite/oapi-sdk-go/core"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var ctxKeyRequestInfo = "------ctxKeyRequestInfo"

type AccessTokenType string

const (
	AccessTokenTypeNone   AccessTokenType = "none_access_token"
	AccessTokenTypeApp    AccessTokenType = "app_access_token"
	AccessTokenTypeTenant AccessTokenType = "tenant_access_token"
	AccessTokenTypeUser   AccessTokenType = "user_access_token"
)

type Opt struct {
	isNotDataField   bool
	pathParams       map[string]interface{}
	queryParams      map[string]interface{}
	userAccessToken  string
	tenantKey        string
	isResponseStream bool
}

type Info struct {
	HttpPath               string                       // request http path
	HttpMethod             string                       // request http method
	QueryParams            string                       // request query
	Input                  interface{}                  // request body
	AccessibleTokenTypeSet map[AccessTokenType]struct{} // request accessible token type
	AccessTokenType        AccessTokenType              // request access token type
	TenantKey              string
	UserAccessToken        string      // user access token
	IsNotDataField         bool        // response body is not data field
	Output                 interface{} // response body data
	Retryable              bool
	optFns                 []OptFn
	IsResponseStream       bool
	IsResponseStreamReal   bool
}

func (i *Info) WithContext(ctx *core.Context) {
	ctx.Set(ctxKeyRequestInfo, i)
}

type OptFn func(*Opt)

func SetUserAccessToken(userAccessToken string) OptFn {
	return func(opt *Opt) {
		opt.userAccessToken = userAccessToken
	}
}

func SetTenantKey(tenantKey string) OptFn {
	return func(opt *Opt) {
		opt.tenantKey = tenantKey
	}
}

func SetQueryParams(queryParams map[string]interface{}) OptFn {
	return func(opt *Opt) {
		opt.queryParams = queryParams
	}
}

func SetPathParams(pathParams map[string]interface{}) OptFn {
	return func(opt *Opt) {
		opt.pathParams = pathParams
	}
}

func SetNotDataField() OptFn {
	return func(opt *Opt) {
		opt.isNotDataField = true
	}
}

func SetResponseStream() OptFn {
	return func(opt *Opt) {
		opt.isResponseStream = true
	}
}

type Request struct {
	*Info
	HTTPRequest         *http.Request
	HTTPResponse        *http.Response
	RequestBody         []byte
	RequestBodyStream   io.Reader
	RequestBodyFilePath string
	ContentType         string
	Err                 error
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s %s", r.HttpMethod, r.url(), r.AccessTokenType)
}

func NewRequestByAuth(httpPath, httpMethod string, input, output interface{}) *Request {
	return &Request{
		Info: &Info{
			HttpPath:        httpPath,
			HttpMethod:      httpMethod,
			AccessTokenType: AccessTokenTypeNone,
			Input:           input,
			Output:          output,
			optFns:          []OptFn{SetNotDataField()},
		},
	}
}

// Deprecated, please use `NewRequestWithNative`
func NewRequest2(httpPath, httpMethod string, accessTokenType AccessTokenType,
	input interface{}, output interface{}, optFns ...OptFn) *Request {
	return NewRequest(httpPath, httpMethod, []AccessTokenType{accessTokenType}, input, output, optFns...)
}

func NewRequestWithNative(httpPath, httpMethod string, accessTokenType AccessTokenType,
	input interface{}, output interface{}, optFns ...OptFn) *Request {
	return NewRequest(httpPath, httpMethod, []AccessTokenType{accessTokenType}, input, output, optFns...)
}

func NewRequest(httpPath, httpMethod string, accessTokenTypes []AccessTokenType,
	input interface{}, output interface{}, optFns ...OptFn) *Request {
	accessibleTokenTypeSet := make(map[AccessTokenType]struct{})
	accessTokenType := accessTokenTypes[0]
	for _, t := range accessTokenTypes {
		if t == AccessTokenTypeTenant {
			accessTokenType = t
		}
		accessibleTokenTypeSet[t] = struct{}{}
	}
	req := &Request{
		Info: &Info{
			HttpPath:               httpPath,
			HttpMethod:             httpMethod,
			AccessTokenType:        accessTokenType,
			AccessibleTokenTypeSet: accessibleTokenTypeSet,
			Input:                  input,
			Output:                 output,
			optFns:                 optFns,
		},
	}
	return req
}

func (r *Request) Init() error {
	opt := &Opt{}
	for _, optFn := range r.optFns {
		optFn(opt)
	}
	r.IsNotDataField = opt.isNotDataField
	r.IsResponseStream = opt.isResponseStream
	if opt.tenantKey != "" {
		if _, ok := r.AccessibleTokenTypeSet[AccessTokenTypeTenant]; ok {
			r.AccessTokenType = AccessTokenTypeTenant
			r.TenantKey = opt.tenantKey
		}
	}
	if opt.userAccessToken != "" {
		if _, ok := r.AccessibleTokenTypeSet[AccessTokenTypeUser]; ok {
			r.AccessTokenType = AccessTokenTypeUser
			r.UserAccessToken = opt.userAccessToken
		}
	}
	if opt.queryParams != nil {
		r.QueryParams = toUrlValues(opt.queryParams).Encode()
	}
	if opt.pathParams != nil {
		httpPath, err := resolvePath(r.HttpPath, opt.pathParams)
		if err != nil {
			return err
		}
		r.HttpPath = httpPath
	}
	return nil
}

func toUrlValues(params map[string]interface{}) url.Values {
	vs := make(url.Values)
	for k, v := range params {
		sv := reflect.ValueOf(v)
		if sv.Kind() == reflect.Slice || sv.Kind() == reflect.Array {
			for i := 0; i < sv.Len(); i++ {
				vs.Add(k, fmt.Sprint(sv.Index(i)))
			}
		} else {
			vs.Set(k, fmt.Sprint(v))
		}
	}
	return vs
}

func resolvePath(path string, pathVar map[string]interface{}) (string, error) {
	tmpPath := path
	newPath := ""
	for {
		i := strings.Index(tmpPath, ":")
		if i == -1 {
			newPath += tmpPath
			break
		}
		newPath += tmpPath[:i]
		subPath := tmpPath[i:]
		j := strings.Index(subPath, "/")
		if j == -1 {
			j = len(subPath)
		}
		varName := subPath[1:j]
		v, ok := pathVar[varName]
		if !ok {
			return "", fmt.Errorf("path:%s, var name:%s not find value", path, varName)
		}
		newPath += fmt.Sprint(v)
		if j == len(subPath) {
			break
		}
		tmpPath = subPath[j:]
	}
	return newPath, nil
}

func (r *Request) url() string {
	path := fmt.Sprintf("/%s/%s", constants.OAPIRootPath, r.HttpPath)
	if r.QueryParams != "" {
		path = fmt.Sprintf("%s?%s", path, r.QueryParams)
	}
	return path
}

func (r *Request) FullUrl(domain string) string {
	return fmt.Sprintf("%s%s", domain, r.url())
}

func (r *Request) DataFilled() bool {
	return r.Output != nil && reflect.ValueOf(r.Output).Elem().IsValid()
}

func GetInfoByCtx(ctx context.Context) *Info {
	return ctx.Value(ctxKeyRequestInfo).(*Info)
}

type UserID struct {
	Type constants.UserIDType
	ID   string
}
