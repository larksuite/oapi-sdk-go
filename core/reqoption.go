package core

import "net/http"

type RequestOption struct {
	TenantKey         string
	UserAccessToken   string
	AppAccessToken    string
	TenantAccessToken string
	NeedHelpDeskAuth  bool
	RequestId         string
	FileUpload        bool
	FileDownload      bool
	Header            http.Header
}

type RequestOptionFunc func(option *RequestOption)

func WithNeedHelpDeskAuth() RequestOptionFunc {
	return func(option *RequestOption) {
		option.NeedHelpDeskAuth = true
	}
}

func WithRequestId(requestId string) RequestOptionFunc {
	return func(option *RequestOption) {
		option.RequestId = requestId
	}
}

func WithTenantKey(tenantKey string) RequestOptionFunc {
	return func(option *RequestOption) {
		option.TenantKey = tenantKey
	}
}

func WithFileUpload() RequestOptionFunc {
	return func(option *RequestOption) {
		option.FileUpload = true
	}
}

func WithFileDownload() RequestOptionFunc {
	return func(option *RequestOption) {
		option.FileDownload = true
	}
}

func WithHeaders(header http.Header) RequestOptionFunc {
	return func(option *RequestOption) {
		option.Header = header
	}
}

func WithUserAccessToken(userAccessToken string) RequestOptionFunc {
	return func(option *RequestOption) {
		option.UserAccessToken = userAccessToken
	}
}

func WithAppAccessToken(appAccessToken string) RequestOptionFunc {
	return func(option *RequestOption) {
		option.AppAccessToken = appAccessToken
	}
}

func WithTenantAccessToken(tenantAccessToken string) RequestOptionFunc {
	return func(option *RequestOption) {
		option.TenantAccessToken = tenantAccessToken
	}
}
