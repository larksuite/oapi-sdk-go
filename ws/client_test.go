package ws

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type wsTenantAccessTokenProviderFunc func(context.Context) (string, error)

func (p wsTenantAccessTokenProviderFunc) RetrieveToken(ctx context.Context) (string, error) {
	return p(ctx)
}

type wsMockClientAssertionProvider struct {
	mu     sync.Mutex
	tokens []*larkcore.Token
	index  int
	calls  int
	auds   []string
}

func (p *wsMockClientAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	p.auds = append(p.auds, aud)
	token := p.tokens[p.index]
	if p.index < len(p.tokens)-1 {
		p.index++
	}
	return token, nil
}

type wsRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f wsRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFetchEndpointWithAppSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != GenEndpointUri {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		if req.AppSecret == "" || req.ClientAssertion != "" {
			t.Fatalf("unexpected bootstrap request: %#v", req)
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer server.Close()

	client := NewClient("app-id", "app-secret", WithDomain(server.URL), WithHttpClient(server.Client()))
	connURL, _, err := client.fetchEndpoint(context.Background())
	if err != nil {
		t.Fatalf("get conn url failed: %v", err)
	}
	if connURL != "wss://example.com/ws" {
		t.Fatalf("unexpected conn url: %s", connURL)
	}
}

func TestFetchEndpointWithCustomHeadersAndUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Tt-Env") != "ppe_test" {
			t.Fatalf("missing X-Tt-Env header: %s", r.Header.Get("X-Tt-Env"))
		}
		if r.Header.Get("X-Use-Ppe") != "1" {
			t.Fatalf("missing X-Use-Ppe header: %s", r.Header.Get("X-Use-Ppe"))
		}
		userAgent := r.Header.Get("User-Agent")
		if !strings.HasPrefix(userAgent, "oapi-sdk-go/") {
			t.Fatalf("unexpected user agent: %s", userAgent)
		}
		if !strings.Contains(userAgent, " source/ws-test") {
			t.Fatalf("missing source in user agent: %s", userAgent)
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer server.Close()

	headers := make(http.Header)
	headers.Set("X-Tt-Env", "ppe_test")
	headers.Set("X-Use-Ppe", "1")

	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithHttpClient(server.Client()),
		WithHeaders(headers),
		WithSource("ws-test"),
	)
	if _, _, err := client.fetchEndpoint(context.Background()); err != nil {
		t.Fatalf("get conn url failed: %v", err)
	}
}

func TestFetchEndpointWithClientAssertionProxy(t *testing.T) {
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy"+GenEndpointUri {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get(larkcore.HeaderXTargetService) != "open.feishu.cn" {
			t.Fatalf("unexpected target service: %s", r.Header.Get(larkcore.HeaderXTargetService))
		}
		var req BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		if req.ClientAssertion != "assertion" || req.AppSecret != "" {
			t.Fatalf("unexpected bootstrap request: %#v", req)
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer proxyServer.Close()

	provider := &wsMockClientAssertionProvider{tokens: []*larkcore.Token{{Value: "assertion", TargetInfo: &larkcore.TargetInfo{TargetService: proxyServer.URL, TargetPrefix: "/proxy"}}}}
	client := NewClient("app-id", "", WithDomain("https://open.feishu.cn"), WithHttpClient(proxyServer.Client()), WithClientAssertionProvider(provider))
	if _, _, err := client.fetchEndpoint(context.Background()); err != nil {
		t.Fatalf("get conn url failed: %v", err)
	}
}

func TestFetchEndpointWithClientAssertionProxyHTTPErrorMsg(t *testing.T) {
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 20050,
			Msg:  "target service unavailable",
		})
	}))
	defer proxyServer.Close()

	provider := &wsMockClientAssertionProvider{tokens: []*larkcore.Token{{Value: "assertion", TargetInfo: &larkcore.TargetInfo{TargetService: proxyServer.URL, TargetPrefix: "/proxy"}}}}
	client := NewClient("app-id", "", WithDomain("https://open.feishu.cn"), WithHttpClient(proxyServer.Client()), WithClientAssertionProvider(provider))
	_, _, err := client.fetchEndpoint(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	serverErr, ok := err.(*ServerError)
	if !ok {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if serverErr.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", serverErr.Code)
	}
	if serverErr.Msg != "target service unavailable" {
		t.Fatalf("unexpected server error message: %q", serverErr.Msg)
	}
}

func TestFetchEndpointRetrieveTokenEachTime(t *testing.T) {
	bodyAssertions := make([]string, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		bodyAssertions = append(bodyAssertions, req.ClientAssertion)
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer server.Close()

	provider := &wsMockClientAssertionProvider{tokens: []*larkcore.Token{{Value: "assertion-1"}, {Value: "assertion-2"}}}
	client := NewClient("app-id", "", WithDomain(server.URL), WithHttpClient(server.Client()), WithClientAssertionProvider(provider))
	if _, _, err := client.fetchEndpoint(context.Background()); err != nil {
		t.Fatalf("first get conn url failed: %v", err)
	}
	if _, _, err := client.fetchEndpoint(context.Background()); err != nil {
		t.Fatalf("second get conn url failed: %v", err)
	}
	if provider.calls != 2 {
		t.Fatalf("expected provider called twice, got %d", provider.calls)
	}
	if len(bodyAssertions) != 2 || bodyAssertions[0] != "assertion-1" || bodyAssertions[1] != "assertion-2" {
		t.Fatalf("unexpected assertions: %#v", bodyAssertions)
	}
}

func TestWithHttpClientFetchesEndpoint(t *testing.T) {
	var requests int
	client := NewClient("app-id", "app-secret",
		WithDomain("https://bootstrap.example"),
		WithHttpClient(&http.Client{Transport: wsRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			if req.URL.String() != "https://bootstrap.example"+GenEndpointUri {
				t.Fatalf("unexpected request URL: %s", req.URL)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"code":0,"data":{"url":"wss://example.com/ws"}}`)),
			}, nil
		})}),
	)

	endpoint, _, err := client.fetchEndpoint(context.Background())
	if err != nil {
		t.Fatalf("fetch endpoint: %v", err)
	}
	if endpoint != "wss://example.com/ws" {
		t.Fatalf("unexpected endpoint: %s", endpoint)
	}
	if requests != 1 {
		t.Fatalf("expected one injected HTTP request, got %d", requests)
	}
}

func TestBuildWSProxyURL(t *testing.T) {
	testCases := map[string]string{
		"proxy.example.com":         "https://proxy.example.com/v1" + GenEndpointUri,
		"https://proxy.example.com": "https://proxy.example.com/v1" + GenEndpointUri,
		"http://proxy.example.com":  "http://proxy.example.com/v1" + GenEndpointUri,
	}

	for targetService, expected := range testCases {
		proxyURL := buildWSProxyURL(targetService, "/v1", GenEndpointUri)
		if proxyURL != expected {
			t.Fatalf("unexpected proxy url for %s: %s", targetService, proxyURL)
		}
	}
}

func TestFetchEndpointWithTenantAccessToken(t *testing.T) {
	for _, secret := range []string{"", "unused-secret"} {
		t.Run(secret, func(t *testing.T) {
			calls := 0
			tokens := []string{"test-tenant-token-1", "test-tenant-token-2"}
			provider := wsTenantAccessTokenProviderFunc(func(context.Context) (string, error) {
				token := tokens[calls]
				calls++
				return token, nil
			})
			requests := 0
			httpClient := &http.Client{Transport: wsRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodPost || req.URL.String() != "https://example.com/callback/ws/v2/endpoint" {
					t.Fatal("TAT bootstrap must use the new endpoint")
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("read bootstrap request: %v", err)
				}
				if string(body) != "{}" {
					t.Fatal("TAT bootstrap must not send application identity or legacy credentials")
				}
				values := req.Header.Values("Authorization")
				if len(values) != 1 || values[0] != "Bearer "+tokens[requests] {
					t.Fatal("Provider must own the unique Authorization header")
				}
				if req.Header.Get("X-Tt-Env") != "ppe_test" || !strings.Contains(req.Header.Get("User-Agent"), " source/tat-test") {
					t.Fatal("TAT bootstrap must preserve routing headers and SDK User-Agent")
				}
				requests++
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"code":0,"data":{"URL":"wss://example.com/ws"}}`))}, nil
			})}
			client := NewClient("unused-app-id", secret,
				WithTenantAccessTokenProvider(provider), WithDomain("https://example.com/"), WithHttpClient(httpClient),
				WithHeaders(http.Header{"Authorization": {"unused-authorization", "another-unused-authorization"}, "X-Tt-Env": {"ppe_test"}}),
				WithSource("tat-test"),
			)
			for i := 0; i < 2; i++ {
				endpoint, config, err := client.fetchEndpoint(context.Background())
				if err != nil || endpoint != "wss://example.com/ws" || config != nil {
					t.Fatalf("fetch TAT endpoint: endpoint=%q, config=%v, err=%v", endpoint, config, err)
				}
			}
			if calls != 2 || requests != 2 {
				t.Fatal("each bootstrap must retrieve a current TAT without SDK caching")
			}
		})
	}
}

func TestFetchEndpointTenantAccessTokenProviderFailures(t *testing.T) {
	providerErr := errors.New("test provider unavailable")
	for _, tc := range []struct {
		name     string
		err      error
		conflict bool
	}{
		{name: "empty token"},
		{name: "provider error", err: providerErr},
		{name: "conflicting providers", conflict: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			provider := wsTenantAccessTokenProviderFunc(func(context.Context) (string, error) {
				calls++
				return "", tc.err
			})
			client := NewClient("app-id", "unused-secret", WithTenantAccessTokenProvider(provider))
			if tc.conflict {
				client.clientAssertionProvider = &wsMockClientAssertionProvider{}
			}
			client.httpClient = &http.Client{Transport: wsRoundTripperFunc(func(*http.Request) (*http.Response, error) {
				t.Fatal("invalid TAT configuration must not issue HTTP requests or fall back to other credentials")
				return nil, errors.New("unexpected request")
			})}
			_, _, err := client.fetchEndpoint(context.Background())
			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatal("Provider error must be preserved")
				}
			} else if _, ok := err.(*ClientError); !ok {
				t.Fatalf("expected configuration ClientError, got %T", err)
			}
			if (tc.conflict && calls != 0) || (!tc.conflict && calls != 1) {
				t.Fatalf("unexpected Provider calls: %d", calls)
			}
		})
	}
}

func TestFetchEndpointTenantAccessTokenUsesExistingErrorHandling(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		body      string
		clientErr bool
		code      int
	}{
		{name: "invalid token", status: 400, body: `{"code":99991663,"msg":"invalid tenant access token"}`, code: 400},
		{name: "unauthorized", status: 401, body: `{}`, code: 401},
		{name: "forbidden", status: 403, body: `{}`, code: 403},
		{name: "missing route", status: 404, body: `{}`, code: 404},
		{name: "rate limited", status: 429, body: `{}`, code: 429},
		{name: "server failure", status: 503, body: `{}`, code: 503},
		{name: "business rejection", status: 200, body: `{"code":19001,"msg":"denied"}`, clientErr: true, code: 19001},
		{name: "busy", status: 200, body: `{"code":1}`, code: SystemBusy},
		{name: "missing endpoint", status: 200, body: `{"code":0}`, code: 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient("", "", WithTenantAccessTokenProvider(wsTenantAccessTokenProviderFunc(func(context.Context) (string, error) {
				return "test-tenant-token", nil
			})))
			client.httpClient = &http.Client{Transport: wsRoundTripperFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body))}, nil
			})}
			_, _, err := client.fetchEndpoint(context.Background())
			if tc.clientErr {
				clientErr, ok := err.(*ClientError)
				if !ok || clientErr.Code != tc.code {
					t.Fatalf("expected terminal client code %d, got %v", tc.code, err)
				}
			} else {
				serverErr, ok := err.(*ServerError)
				if !ok || serverErr.Code != tc.code {
					t.Fatalf("expected retryable server code %d, got %v", tc.code, err)
				}
			}
		})
	}
}

func TestTenantAccessTokenUsesExistingReconnectConfig(t *testing.T) {
	client := NewClient("", "", WithTenantAccessTokenProvider(wsTenantAccessTokenProviderFunc(func(context.Context) (string, error) {
		return "test-tenant-token", nil
	})))
	if client.reconnectCount != -1 {
		t.Fatal("TAT must preserve the existing default reconnect count")
	}
	for _, tc := range []struct{ configured, want int }{{-1, -1}, {0, 0}, {2, 2}, {10, 10}} {
		client.configure(&ClientConfig{ReconnectCount: tc.configured})
		if got, _, _, _ := client.configSnapshot(); got != tc.want {
			t.Fatalf("TAT ReconnectCount=%d, want %d", got, tc.want)
		}
	}
	legacy := NewClient("app-id", "secret", WithTenantAccessTokenProvider(nil))
	legacy.configure(&ClientConfig{ReconnectCount: -1})
	if legacy.reconnectCount != -1 {
		t.Fatal("nil TAT Provider must preserve legacy retry configuration")
	}
}
