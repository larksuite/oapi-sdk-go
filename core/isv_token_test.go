package larkcore

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type isvTokenCache struct {
	localCache
	writes int
	ttl    time.Duration
	err    error
}

func (c *isvTokenCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	c.writes++
	c.ttl = ttl
	if c.err != nil {
		return c.err
	}
	return c.localCache.Set(ctx, key, value, ttl)
}

func newISVTokenConfig(server *httptest.Server) *Config {
	config := mockConfig()
	config.AppId = "cli_isv"
	config.AppType = AppTypeMarketplace
	config.AppSecret = "unused-secret"
	config.EnableTokenCache = true
	config.OAuthBaseUrl = server.URL
	config.HttpClient = server.Client()
	config.ReqTimeout = time.Second
	config.TokenCache = &isvTokenCache{}
	config.ClientAssertionProvider = &mockClientAssertionProvider{token: &Token{Value: "test-assertion"}}
	return config
}

func writeISVTokenResponse(t *testing.T, w http.ResponseWriter, response OAuthTokenResp) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Errorf("encode test token response: %v", err)
	}
}

func TestISVTenantTokenExchangeAndCache(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != OAuthTokenUrlPath || r.Method != http.MethodPost {
			t.Error("unexpected token endpoint or method")
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			t.Error("expected JSON OAuth request")
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode token body: %v", err)
			return
		}
		if body["grant_type"] != GrantTypeClientCredentials ||
			body["client_id"] != "cli_isv" ||
			body["client_assertion_type"] != ClientAssertionTypeJWTBearer ||
			body["client_assertion"] != "test-assertion" {
			t.Error("incorrect ISV client assertion parameters")
		}
		for _, field := range []string{"client_secret", "app_secret", "app_ticket", "app_access_token"} {
			if _, ok := body[field]; ok {
				t.Errorf("unexpected credential field: %s", field)
			}
		}
		writeISVTokenResponse(t, w, OAuthTokenResp{
			AccessToken: "test-token-" + body["tenant_key"],
			TokenType:   "Bearer", ExpiresIn: 7200,
		})
	}))
	defer server.Close()

	config := newISVTokenConfig(server)
	manager := TokenManager{cache: config.TokenCache}
	for _, tenant := range []string{"tenant-A", "tenant-B", "tenant-A", "tenant-B"} {
		token, err := manager.getTenantAccessToken(context.Background(), config, tenant, "")
		if err != nil {
			t.Fatalf("exchange ISV TAT: %v", err)
		}
		if token != "test-token-"+tenant {
			t.Fatal("token belongs to the wrong tenant")
		}
	}
	if calls != 2 || config.ClientAssertionProvider.(*mockClientAssertionProvider).calls != 2 {
		t.Fatal("each tenant should exchange once and then use its cache")
	}
	tokenCache := config.TokenCache.(*isvTokenCache)
	if tokenCache.writes != 2 || tokenCache.ttl != 117*time.Minute {
		t.Fatal("unexpected cache writes or TTL")
	}
}

func TestISVTenantTokenWithoutAdditionalTimeoutConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeISVTokenResponse(t, w, OAuthTokenResp{AccessToken: "test-token", TokenType: "Bearer", ExpiresIn: 7200})
	}))
	defer server.Close()
	config := newISVTokenConfig(server)
	config.ReqTimeout = 0
	manager := TokenManager{cache: config.TokenCache}
	if _, err := manager.getTenantAccessToken(context.Background(), config, "tenant-A", ""); err != nil {
		t.Fatalf("ISV mode must not impose a new timeout configuration requirement: %v", err)
	}
}

func TestClientAssertionProviderErrorUsesExistingCode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	providerErr := errors.New("provider unavailable")
	for _, appType := range []AppType{AppTypeSelfBuilt, AppTypeMarketplace} {
		config := mockConfig()
		config.AppType = appType
		config.OAuthBaseUrl = "https://accounts.feishu.cn"
		config.EnableTokenCache = false
		config.ClientAssertionProvider = &mockClientAssertionProvider{err: providerErr}
		manager := TokenManager{}
		_, err := manager.getTenantTokenByClientAssertion(ctx, config, "tenant-A")
		var codeErr *CodeError
		if !errors.As(err, &codeErr) || codeErr.Code != ErrCodeClientAssertionRetrieveFailed || codeErr.Msg != providerErr.Error() {
			t.Fatalf("app type %v did not preserve the existing provider error", appType)
		}
	}
}

func TestISVTenantTokenRejectsFailedResponses(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   OAuthTokenResp
	}{
		{name: "http failure", status: http.StatusBadRequest, body: OAuthTokenResp{AccessToken: "test-token", ExpiresIn: 7200, TokenType: "Bearer"}},
		{name: "business failure", status: http.StatusOK, body: OAuthTokenResp{Code: 20140, AccessToken: "test-token", ExpiresIn: 7200, TokenType: "Bearer"}},
		{name: "oauth failure", status: http.StatusOK, body: OAuthTokenResp{Error: "invalid_client", AccessToken: "test-token", ExpiresIn: 7200, TokenType: "Bearer"}},
		{name: "empty token", status: http.StatusOK, body: OAuthTokenResp{ExpiresIn: 7200, TokenType: "Bearer"}},
		{name: "invalid expiration", status: http.StatusOK, body: OAuthTokenResp{AccessToken: "test-token", TokenType: "Bearer"}},
		{name: "invalid token type", status: http.StatusOK, body: OAuthTokenResp{AccessToken: "test-token", ExpiresIn: 7200, TokenType: "Other"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.status)
				writeISVTokenResponse(t, w, test.body)
			}))
			defer server.Close()
			config := newISVTokenConfig(server)
			logger := &recordingLogger{}
			config.Logger = logger
			manager := TokenManager{cache: config.TokenCache}
			if _, err := manager.getTenantAccessToken(context.Background(), config, "tenant-A", ""); err == nil {
				t.Fatal("expected token exchange failure")
			}
			if config.TokenCache.(*isvTokenCache).writes != 0 {
				t.Fatal("failed token response was cached")
			}
			if !strings.Contains(logger.String(), "ISV tenant token") {
				t.Fatal("token exchange failure was not logged")
			}
			for _, credential := range []string{"test-token", "test-assertion", "unused-secret"} {
				if strings.Contains(logger.String(), credential) {
					t.Fatal("token exchange failure log leaked a credential")
				}
			}
		})
	}
}

func TestISVTenantTokenShortLifetimeAndCacheFailure(t *testing.T) {
	for _, test := range []struct {
		name      string
		expiresIn int
		cacheErr  error
	}{
		{name: "short lifetime", expiresIn: 60},
		{name: "cache failure", expiresIn: 7200, cacheErr: errors.New("cache unavailable")},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeISVTokenResponse(t, w, OAuthTokenResp{AccessToken: "test-token", TokenType: "Bearer", ExpiresIn: test.expiresIn})
			}))
			defer server.Close()
			config := newISVTokenConfig(server)
			config.TokenCache.(*isvTokenCache).err = test.cacheErr
			logger := &recordingLogger{}
			config.Logger = logger
			manager := TokenManager{cache: config.TokenCache}
			token, err := manager.getTenantAccessToken(context.Background(), config, "tenant-A", "")
			if err != nil {
				t.Fatalf("exchange token: %v", err)
			}
			if token != "test-token" {
				t.Fatal("cache handling must preserve the exchanged token")
			}
			tokenCache := config.TokenCache.(*isvTokenCache)
			if tokenCache.writes != 1 {
				t.Fatal("expected one cache write attempt")
			}
			if test.expiresIn == 60 && tokenCache.ttl != 0 {
				t.Fatal("short-lived token must use the existing zero TTL behavior")
			}
			if test.cacheErr != nil && !strings.Contains(logger.String(), "client assertion tenantAccessToken save cache") {
				t.Fatal("cache write failure must be logged")
			}
		})
	}
}

func TestISVTenantTokenRequiresTenantBeforeProvider(t *testing.T) {
	config := mockConfig()
	config.AppType = AppTypeMarketplace
	config.BaseUrl = "https://open.feishu.cn"
	config.ReqTimeout = time.Second
	provider := &mockClientAssertionProvider{token: &Token{Value: "test-assertion"}}
	config.ClientAssertionProvider = provider
	manager := TokenManager{cache: config.TokenCache}
	for _, tenant := range []string{"", "   "} {
		if _, err := manager.getTenantAccessToken(context.Background(), config, tenant, ""); err == nil {
			t.Fatal("expected missing tenant error")
		}
	}
	if provider.calls != 0 {
		t.Fatal("invalid tenant must fail before calling provider")
	}
}

func TestISVTenantTokenProxy(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy"+OAuthTokenUrlPath || r.Header.Get(HeaderXTargetService) != "accounts.feishu.cn" {
			t.Error("incorrect OAuth proxy routing")
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode proxy body: %v", err)
			return
		}
		if body["tenant_key"] != "tenant-A" {
			t.Error("proxy request lost target tenant")
		}
		writeISVTokenResponse(t, w, OAuthTokenResp{AccessToken: "test-token", TokenType: "Bearer", ExpiresIn: 7200})
	}))
	defer server.Close()
	config := newISVTokenConfig(server)
	config.OAuthBaseUrl = "https://accounts.feishu.cn"
	provider := config.ClientAssertionProvider.(*mockClientAssertionProvider)
	provider.token.TargetInfo = &TargetInfo{TargetService: server.URL, TargetPrefix: "/proxy"}
	logger := &recordingLogger{}
	config.Logger = logger
	config.LogReqAtDebug = true
	manager := TokenManager{cache: config.TokenCache}
	if _, err := manager.getTenantAccessToken(context.Background(), config, "tenant-A", ""); err != nil {
		t.Fatalf("exchange through proxy: %v", err)
	}
	for _, value := range []string{"test-assertion", "test-token", "unused-secret"} {
		if strings.Contains(logger.String(), value) {
			t.Fatal("OAuth debug logs leaked a credential")
		}
	}
	if len(provider.auds) != 1 || provider.auds[0] != "accounts.feishu.cn" {
		t.Fatal("provider audience must identify the OAuth host")
	}
}

func TestISVTenantTokenCacheSeparatesAppAndEndpoint(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode token body: %v", err)
			return
		}
		writeISVTokenResponse(t, w, OAuthTokenResp{
			AccessToken: strings.Repeat("x", 4096) + body["client_id"] + r.URL.Path,
			TokenType:   "Bearer", ExpiresIn: 7200,
		})
	}))
	defer server.Close()
	sharedCache := &isvTokenCache{}
	manager := TokenManager{cache: sharedCache}
	for i := 0; i < 2; i++ {
		for _, test := range []struct{ appID, prefix string }{
			{appID: "cli-A", prefix: "/region-A"},
			{appID: "cli-B", prefix: "/region-A"},
			{appID: "cli-A", prefix: "/region-B"},
		} {
			config := newISVTokenConfig(server)
			config.AppId = test.appID
			config.OAuthBaseUrl = server.URL + test.prefix
			config.TokenCache = sharedCache
			token, err := manager.getTenantAccessToken(context.Background(), config, "tenant-A", "")
			if err != nil {
				t.Fatalf("exchange with shared cache: %v", err)
			}
			if token != strings.Repeat("x", 4096)+test.appID+test.prefix+OAuthTokenUrlPath {
				t.Fatal("cached token was truncated or belongs to another application or endpoint")
			}
		}
	}
	if calls != 3 {
		t.Fatal("applications and endpoints must have independent cache entries")
	}
}
