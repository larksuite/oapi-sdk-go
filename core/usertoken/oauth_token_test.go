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

package usertoken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type mockClientAssertionProvider struct {
	token *larkcore.Token
	err   error
	auds  []string
}

func (p *mockClientAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	p.auds = append(p.auds, aud)
	return p.token, p.err
}

func newTestConfig(server *httptest.Server, provider larkcore.ClientAssertionProvider) *larkcore.Config {
	return &larkcore.Config{
		BaseUrl:                 "https://open.feishu.cn",
		OAuthBaseUrl:            server.URL,
		AppId:                   "cli_a",
		AppSecret:               "",
		ClientAssertionProvider: provider,
		EnableTokenCache:        true,
		HttpClient:              server.Client(),
		Logger:                  larkcore.NewDefaultLogger(larkcore.LogLevelInfo),
		Serializable:            &larkcore.DefaultSerialization{},
	}
}

func newAppSecretTestConfig(server *httptest.Server) *larkcore.Config {
	config := newTestConfig(server, nil)
	config.AppSecret = "app-secret"
	return config
}

func TestOAuthTokenCreate(t *testing.T) {
	provider := &mockClientAssertionProvider{token: &larkcore.Token{Value: "client-assertion"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != larkcore.OAuthTokenUrlPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req oauthTokenRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode body failed: %v", err)
		}
		if req.GrantType != larkcore.GrantTypeAuthorizationCode {
			t.Fatalf("unexpected grant type: %s", req.GrantType)
		}
		if req.ClientAssertion != "client-assertion" || req.ClientID != "cli_a" {
			t.Fatalf("unexpected client assertion request: %#v", req)
		}
		if req.Code != "code" || req.RedirectUri != "https://example.com/cb" || req.CodeVerifier != "verifier" {
			t.Fatalf("unexpected request body: %#v", req)
		}
		_ = json.NewEncoder(w).Encode(&oauthTokenResponseBody{
			AccessToken:  "user-token",
			TokenType:    "Bearer",
			ExpiresIn:    7200,
			RefreshToken: "refresh-token",
			Scope:        "contact:user.base:readonly",
		})
	}))
	defer server.Close()

	oauthToken := NewOAuthToken(newTestConfig(server, provider))
	resp, err := oauthToken.Create(context.Background(), NewCreateOAuthTokenReqBuilder().
		Code("code").
		RedirectUri("https://example.com/cb").
		CodeVerifier("verifier").
		Build())
	if err != nil {
		t.Fatalf("create oauth token failed: %v", err)
	}
	if !resp.Success() || larkcore.StringValue(resp.Data.AccessToken) != "user-token" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if len(provider.auds) != 1 || provider.auds[0] != server.URL {
		t.Fatalf("unexpected auds: %#v", provider.auds)
	}
}

func TestOAuthTokenRefresh(t *testing.T) {
	provider := &mockClientAssertionProvider{token: &larkcore.Token{Value: "client-assertion"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req oauthTokenRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode body failed: %v", err)
		}
		if req.GrantType != larkcore.GrantTypeRefreshToken {
			t.Fatalf("unexpected grant type: %s", req.GrantType)
		}
		if req.RefreshToken != "refresh-token" {
			t.Fatalf("unexpected refresh token: %s", req.RefreshToken)
		}
		_ = json.NewEncoder(w).Encode(&oauthTokenResponseBody{AccessToken: "user-token", ExpiresIn: 7200})
	}))
	defer server.Close()

	oauthToken := NewOAuthToken(newTestConfig(server, provider))
	resp, err := oauthToken.Refresh(context.Background(), NewRefreshOAuthTokenReqBuilder().RefreshToken("refresh-token").Build())
	if err != nil {
		t.Fatalf("refresh oauth token failed: %v", err)
	}
	if larkcore.StringValue(resp.Data.AccessToken) != "user-token" {
		t.Fatalf("unexpected access token: %#v", resp.Data.AccessToken)
	}
}

func TestOAuthTokenCreateWithAppSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != larkcore.OAuthTokenUrlPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req oauthTokenRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode body failed: %v", err)
		}
		if req.ClientID != "cli_a" || req.ClientSecret != "app-secret" {
			t.Fatalf("unexpected app secret request: %#v", req)
		}
		if req.ClientAssertion != "" || req.ClientAssertionType != "" {
			t.Fatalf("unexpected client assertion fields: %#v", req)
		}
		if req.GrantType != larkcore.GrantTypeAuthorizationCode || req.Code != "code" {
			t.Fatalf("unexpected request body: %#v", req)
		}
		_ = json.NewEncoder(w).Encode(&oauthTokenResponseBody{AccessToken: "user-token", ExpiresIn: 7200})
	}))
	defer server.Close()

	oauthToken := NewOAuthToken(newAppSecretTestConfig(server))
	resp, err := oauthToken.Create(context.Background(), NewCreateOAuthTokenReqBuilder().Code("code").Build())
	if err != nil {
		t.Fatalf("create oauth token with app secret failed: %v", err)
	}
	if larkcore.StringValue(resp.Data.AccessToken) != "user-token" {
		t.Fatalf("unexpected access token: %#v", resp.Data.AccessToken)
	}
}

func TestOAuthTokenRefreshWithAppSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req oauthTokenRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode body failed: %v", err)
		}
		if req.ClientID != "cli_a" || req.ClientSecret != "app-secret" {
			t.Fatalf("unexpected app secret request: %#v", req)
		}
		if req.GrantType != larkcore.GrantTypeRefreshToken || req.RefreshToken != "refresh-token" {
			t.Fatalf("unexpected request body: %#v", req)
		}
		_ = json.NewEncoder(w).Encode(&oauthTokenResponseBody{AccessToken: "user-token", ExpiresIn: 7200})
	}))
	defer server.Close()

	oauthToken := NewOAuthToken(newAppSecretTestConfig(server))
	resp, err := oauthToken.Refresh(context.Background(), NewRefreshOAuthTokenReqBuilder().RefreshToken("refresh-token").Build())
	if err != nil {
		t.Fatalf("refresh oauth token with app secret failed: %v", err)
	}
	if larkcore.StringValue(resp.Data.AccessToken) != "user-token" {
		t.Fatalf("unexpected access token: %#v", resp.Data.AccessToken)
	}
}

func TestOAuthTokenRejectsMissingCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be called")
	}))
	defer server.Close()

	config := newTestConfig(server, nil)
	oauthToken := NewOAuthToken(config)
	_, err := oauthToken.Create(context.Background(), NewCreateOAuthTokenReqBuilder().Code("code").Build())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*larkcore.CodeError)
	if !ok || codeErr.Code != larkcore.ErrCodeAppSecretAndClientAssertionEmpty {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestOAuthTokenReturnsOAuthErrorForNonOK(t *testing.T) {
	provider := &mockClientAssertionProvider{token: &larkcore.Token{Value: "client-assertion"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&oauthTokenResponseBody{
			Code:             20001,
			Error:            "invalid_client",
			ErrorDescription: "client assertion invalid",
		})
	}))
	defer server.Close()

	oauthToken := NewOAuthToken(newTestConfig(server, provider))
	_, err := oauthToken.Create(context.Background(), NewCreateOAuthTokenReqBuilder().Code("code").Build())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	oauthErr, ok := err.(*OAuthError)
	if !ok {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if oauthErr.Code != 20001 || oauthErr.ErrorDescription != "client assertion invalid" {
		t.Fatalf("unexpected oauth error: %#v", oauthErr)
	}
}

func TestOAuthTokenProxyKeepsCustomHeaders(t *testing.T) {
	provider := &mockClientAssertionProvider{token: &larkcore.Token{Value: "client-assertion"}}
	proxyServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy"+larkcore.OAuthTokenUrlPath {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get(larkcore.HeaderXTargetService) == "" {
			t.Fatalf("missing target service header")
		}
		if r.Header.Get("X-Custom") != "custom-value" {
			t.Fatalf("missing custom header")
		}
		_ = json.NewEncoder(w).Encode(&oauthTokenResponseBody{AccessToken: "user-token", ExpiresIn: 7200})
	}))
	defer proxyServer.Close()
	provider.token.TargetInfo = &larkcore.TargetInfo{TargetService: proxyServer.Listener.Addr().String(), TargetPrefix: "/proxy"}

	config := newTestConfig(proxyServer, provider)
	oauthToken := NewOAuthToken(config)
	headers := make(http.Header)
	headers.Set("X-Custom", "custom-value")
	_, err := oauthToken.Create(context.Background(), NewCreateOAuthTokenReqBuilder().Code("code").Build(), larkcore.WithHeaders(headers))
	if err != nil {
		t.Fatalf("create oauth token with proxy failed: %v", err)
	}
}
