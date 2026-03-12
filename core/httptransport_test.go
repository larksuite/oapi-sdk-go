/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkcore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
)

func TestSendPost(t *testing.T) {
	config := mockConfig()
	_, err := Request(context.Background(), &ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    "/",
		Body: map[string]interface{}{
			"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		},
		SupportedAccessTokenTypes: []AccessTokenType{AccessTokenTypeUser},
	}, config, WithUserAccessToken("key"))

	if err != nil {
		t.Errorf("TestSendPost failed ,%v", err)
		return
	}
	fmt.Println("ok")

}

type closeTrackingBody struct {
	closed int32
}

func (b *closeTrackingBody) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (b *closeTrackingBody) Close() error {
	atomic.StoreInt32(&b.closed, 1)
	return nil
}

type httpClientStub struct {
	resp *http.Response
	err  error
}

func (c httpClientStub) Do(req *http.Request) (*http.Response, error) {
	return c.resp, c.err
}

type noopLogger struct{}

func (noopLogger) Debug(context.Context, ...interface{}) {}
func (noopLogger) Info(context.Context, ...interface{})  {}
func (noopLogger) Warn(context.Context, ...interface{})  {}
func (noopLogger) Error(context.Context, ...interface{}) {}

func TestDoSend_CloseBodyOnGatewayTimeout(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	body := &closeTrackingBody{}
	client := httpClientStub{resp: &http.Response{
		StatusCode: http.StatusGatewayTimeout,
		Header:     make(http.Header),
		Body:       body,
		Request:    req,
	}}

	resp, err := doSend(context.Background(), req, client, noopLogger{})
	if resp != nil {
		t.Fatalf("expect nil resp, got: %#v", resp)
	}
	if err == nil {
		t.Fatalf("expect error, got nil")
	}
	if _, ok := err.(*ServerTimeoutError); !ok {
		t.Fatalf("expect *ServerTimeoutError, got: %T (%v)", err, err)
	}
	if atomic.LoadInt32(&body.closed) != 1 {
		t.Fatalf("expect resp body closed on gateway timeout")
	}
}
