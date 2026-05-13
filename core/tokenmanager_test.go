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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockClientAssertionProvider struct {
	token *Token
	err   error
	calls int
	auds  []string
}

func (p *mockClientAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*Token, error) {
	p.calls++
	p.auds = append(p.auds, aud)
	return p.token, p.err
}

func TestTokenManagerSetAndGet(t *testing.T) {
	config := mockConfig()
	cache := &localCache{}
	tokenManager := TokenManager{cache: cache}

	err := tokenManager.set(context.Background(), tenantAccessTokenKey(config.AppId, "tenantKey"), "tokenValue", time.Minute)
	if err != nil {
		t.Errorf("set key failed ,%v", err)
	}

	token, err := tokenManager.getTenantAccessToken(context.Background(), config, "tenantKey", "")
	if err != nil {
		t.Errorf("get key failed ,%v", err)
	}

	if token == "" {
		t.Errorf("get key failed ,%v", err)
	}

}

func TestGetTenantAccessTokenByClientAssertion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != OAuthTokenUrlPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req OAuthTokenReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode body failed: %v", err)
		}
		if req.ClientAssertion != "client-assertion" {
			t.Fatalf("unexpected assertion: %s", req.ClientAssertion)
		}
		if req.ClientID != "cli_a" {
			t.Fatalf("unexpected client id: %s", req.ClientID)
		}
		_ = json.NewEncoder(w).Encode(&OAuthTokenResp{AccessToken: "tenant-token", ExpiresIn: 7200})
	}))
	defer server.Close()

	config := mockConfig()
	config.BaseUrl = server.URL
	config.AppId = "cli_a"
	config.EnableTokenCache = true
	config.HttpClient = server.Client()
	config.ClientAssertionProvider = &mockClientAssertionProvider{token: &Token{Value: "client-assertion"}}

	manager := TokenManager{cache: &localCache{}}
	token, err := manager.getTenantAccessToken(context.Background(), config, "tenantKey", "")
	if err != nil {
		t.Fatalf("get tenant access token failed: %v", err)
	}
	if token != "tenant-token" {
		t.Fatalf("unexpected token: %s", token)
	}
}

func TestGetTenantAccessTokenByClientAssertionWithProxy(t *testing.T) {
	proxyServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy"+OAuthTokenUrlPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get(HeaderXTargetService) == "" {
			t.Fatalf("missing target service header")
		}
		_ = json.NewEncoder(w).Encode(&OAuthTokenResp{AccessToken: "tenant-token", ExpiresIn: 7200})
	}))
	defer proxyServer.Close()

	config := mockConfig()
	config.BaseUrl = "https://open.feishu.cn"
	config.EnableTokenCache = true
	config.HttpClient = proxyServer.Client()
	config.ClientAssertionProvider = &mockClientAssertionProvider{token: &Token{Value: "client-assertion", TargetInfo: &TargetInfo{TargetService: proxyServer.Listener.Addr().String(), TargetPrefix: "/proxy"}}}

	manager := TokenManager{cache: &localCache{}}
	token, err := manager.getTenantAccessToken(context.Background(), config, "tenantKey", "")
	if err != nil {
		t.Fatalf("get tenant access token failed: %v", err)
	}
	if token != "tenant-token" {
		t.Fatalf("unexpected token: %s", token)
	}
}

func TestGetTenantAccessTokenByClientAssertionEmptyToken(t *testing.T) {
	config := mockConfig()
	config.ClientAssertionProvider = &mockClientAssertionProvider{token: &Token{}}

	manager := TokenManager{cache: &localCache{}}
	_, err := manager.getTenantAccessToken(context.Background(), config, "tenantKey", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*CodeError)
	if !ok || codeErr.Code != ErrCodeClientAssertionTokenEmpty {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestGetTenantAccessTokenByClientAssertionRetrieveFailed(t *testing.T) {
	config := mockConfig()
	config.ClientAssertionProvider = &mockClientAssertionProvider{err: ErrAppTicketIsEmpty}

	manager := TokenManager{cache: &localCache{}}
	_, err := manager.getTenantAccessToken(context.Background(), config, "tenantKey", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*CodeError)
	if !ok || codeErr.Code != ErrCodeClientAssertionRetrieveFailed {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestGetTenantAccessTokenByClientAssertionWithoutCache(t *testing.T) {
	provider := &mockClientAssertionProvider{token: &Token{Value: "client-assertion"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(&OAuthTokenResp{AccessToken: "tenant-token", ExpiresIn: 7200})
	}))
	defer server.Close()

	config := mockConfig()
	config.BaseUrl = server.URL
	config.EnableTokenCache = false
	config.HttpClient = server.Client()
	config.ClientAssertionProvider = provider

	manager := TokenManager{cache: &localCache{}}
	if _, err := manager.getTenantAccessToken(context.Background(), config, "tenantKey", ""); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if _, err := manager.getTenantAccessToken(context.Background(), config, "tenantKey", ""); err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if provider.calls != 2 {
		t.Fatalf("expected provider called twice, got %d", provider.calls)
	}
}

func TestGetTenantAccessTokenByClientAssertionOAuthErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(&OAuthTokenResp{
			Code:             20050,
			Error:            "server_error",
			ErrorDescription: "An unexpected server error occurred. Please retry your request.",
		})
	}))
	defer server.Close()

	config := mockConfig()
	config.BaseUrl = server.URL
	config.EnableTokenCache = true
	config.HttpClient = server.Client()
	config.ClientAssertionProvider = &mockClientAssertionProvider{token: &Token{Value: "client-assertion"}}

	manager := TokenManager{cache: &localCache{}}
	_, err := manager.getTenantAccessToken(context.Background(), config, "tenantKey", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	codeErr, ok := err.(*CodeError)
	if !ok {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if codeErr.Code != 20050 {
		t.Fatalf("unexpected error code: %d", codeErr.Code)
	}
	if codeErr.Msg != "An unexpected server error occurred. Please retry your request." {
		t.Fatalf("unexpected error message: %s", codeErr.Msg)
	}
}

func TestGetTenantAccessTokenByClientAssertionManualCallWithoutCache(t *testing.T) {
	provider := &mockClientAssertionProvider{token: &Token{Value: "client-assertion"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(&OAuthTokenResp{AccessToken: "tenant-token", ExpiresIn: 7200})
	}))
	defer server.Close()

	config := mockConfig()
	config.BaseUrl = server.URL
	config.EnableTokenCache = false
	config.HttpClient = server.Client()
	config.ClientAssertionProvider = provider

	manager := TokenManager{cache: &localCache{}}
	token, err := manager.getTenantTokenByClientAssertion(context.Background(), config, "tenantKey")
	if err != nil {
		t.Fatalf("manual get tenant token failed: %v", err)
	}
	if token != "tenant-token" {
		t.Fatalf("unexpected token: %s", token)
	}
	if provider.calls != 1 {
		t.Fatalf("expected provider called once, got %d", provider.calls)
	}
}

func TestGetAppAccessTokenBlockedInClientAssertionMode(t *testing.T) {
	config := mockConfig()
	config.ClientAssertionProvider = &mockClientAssertionProvider{token: &Token{Value: "client-assertion"}}

	manager := TokenManager{cache: &localCache{}}
	_, err := manager.getAppAccessToken(context.Background(), config, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*CodeError)
	if !ok || codeErr.Code != ErrCodeClientAssertionProviderNotConfigured {
		t.Fatalf("unexpected error: %#v", err)
	}
}
