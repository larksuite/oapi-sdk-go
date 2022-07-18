package larkcore

import "net/url"

type HttpReq struct {
	HttpMethod                string
	ApiPath                   string
	Body                      interface{}
	QueryParams               QueryParams
	PathParams                PathParams
	SupportedAccessTokenTypes []AccessTokenType
}

type PathParams map[string][]string

func (u PathParams) Get(key string) string {
	vs := u[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}
func (u PathParams) Set(key, value string) {
	u[key] = []string{value}
}

func (u PathParams) Encode() string {
	return url.Values(u).Encode()
}

type QueryParams map[string][]string

func (u QueryParams) Get(key string) string {
	vs := u[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}
func (u QueryParams) Set(key, value string) {
	u[key] = []string{value}
}

func (u QueryParams) Encode() string {
	return url.Values(u).Encode()
}
