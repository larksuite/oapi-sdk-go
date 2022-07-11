package larkcore

import (
	"net/http"
)

// HttpClient :sdk-core面向HttpClient接口编程，实现core与httpclient解耦
//1.可以适配所有基于go-sdk内置httpclient构建的三方httpclient
//2.可以方便的对httpclient进行mock, 方便编写单元测试
type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}
