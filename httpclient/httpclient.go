package httpclient

import (
	"net/http"
	"time"
)

func NewHttpClient(reqTimeout time.Duration) HttpClient {
	if reqTimeout == 0 {
		return http.DefaultClient
	}
	return &http.Client{Timeout: reqTimeout}
}

// HttpClient :sdk-core面向HttpClient接口编程，实现core与httpclient解耦
//1.可以适配所有基于go-sdk内置httpclient构建的三方httpclient
//2.可以方便的对httpclient进行mock, 方便编写单元测试
type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}
