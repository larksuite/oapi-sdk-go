package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/accesstoken/authorizationcode"
	"github.com/larksuite/oapi-sdk-go/v3/core/accesstoken/refreshtoken"
)

type isvAssertionProvider struct {
	calls int64
}

func (p *isvAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	if _, ok := ctx.Deadline(); !ok {
		return nil, fmt.Errorf("assertion provider requires a deadline")
	}
	atomic.AddInt64(&p.calls, 1)
	return &larkcore.Token{Value: "test-assertion"}, nil
}

func TestISVClientAssertionLifecycle(t *testing.T) {
	var tokenCalls int64
	var businessCalls int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case larkcore.OAuthTokenUrlPath:
			atomic.AddInt64(&tokenCalls, 1)
			var body map[string]string
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("parse OAuth request: %v", err)
				return
			}
			token := "test-user-token"
			if body["grant_type"] == larkcore.GrantTypeClientCredentials {
				token = "test-tenant-token-" + body["tenant_key"]
			} else if _, ok := body["tenant_key"]; ok {
				t.Error("UAT exchange must not include a tenant selector")
			}
			if err := json.NewEncoder(w).Encode(larkcore.OAuthTokenResp{
				AccessToken: token, TokenType: "Bearer", ExpiresIn: 7200,
			}); err != nil {
				t.Errorf("encode OAuth response: %v", err)
			}
		case "/resource":
			atomic.AddInt64(&businessCalls, 1)
			if r.Header.Get("Authorization") != "Bearer test-tenant-token-tenant-A" &&
				r.Header.Get("Authorization") != "Bearer test-user-token" {
				t.Error("business request did not use the expected access token")
			}
			if err := json.NewEncoder(w).Encode(larkcore.CodeError{Code: 0}); err != nil {
				t.Errorf("encode business response: %v", err)
			}
		default:
			t.Error("unexpected request: keyless ISV must not request an app ticket or app token")
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	provider := &isvAssertionProvider{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := NewClient(t.Name(), "", WithMarketplaceApp(), WithClientAssertionProvider(provider),
		WithOpenBaseUrl(server.URL), WithOAuthBaseUrl(server.URL),
		WithHttpClient(server.Client()), WithReqTimeout(time.Second))
	if atomic.LoadInt64(&tokenCalls) != 0 || atomic.LoadInt64(&provider.calls) != 0 {
		t.Fatal("initializing the keyless client must not request credentials")
	}
	for i := 0; i < 2; i++ {
		if _, err := client.Get(ctx, "/resource", nil, larkcore.AccessTokenTypeTenant, larkcore.WithTenantKey("tenant-A")); err != nil {
			t.Fatalf("business API with automatic TAT: %v", err)
		}
	}
	if atomic.LoadInt64(&tokenCalls) != 1 {
		t.Fatal("second business request must reuse the tenant token")
	}

	uat, err := client.AccessToken.RetrieveByAuthorizationCode(ctx, authorizationcode.NewTokenRequestBuilder().Code("test-code").Build())
	if err != nil {
		t.Fatalf("ISV authorization code exchange: %v", err)
	}
	if _, err := client.AccessToken.Refresh(ctx, refreshtoken.NewTokenRequestBuilder().RefreshToken("test-refresh").Build()); err != nil {
		t.Fatalf("ISV refresh exchange: %v", err)
	}
	if _, err := client.Get(context.Background(), "/resource", nil, larkcore.AccessTokenTypeUser,
		larkcore.WithUserAccessToken(larkcore.StringValue(uat.Data.AccessToken))); err != nil {
		t.Fatalf("business API with explicit UAT: %v", err)
	}
	if atomic.LoadInt64(&tokenCalls) != 3 || atomic.LoadInt64(&provider.calls) != 3 || atomic.LoadInt64(&businessCalls) != 3 {
		t.Fatal("unexpected exchanges: explicit UAT must not trigger automatic TAT retrieval")
	}
	manualClient := NewClient(t.Name(), "", WithMarketplaceApp(), WithClientAssertionProvider(provider),
		WithOpenBaseUrl(server.URL), WithOAuthBaseUrl(server.URL), WithEnableTokenCache(false),
		WithHttpClient(server.Client()), WithReqTimeout(time.Second))
	if _, err := manualClient.Get(context.Background(), "/resource", nil, larkcore.AccessTokenTypeTenant,
		larkcore.WithTenantKey("tenant-A")); err == nil {
		t.Fatal("disabled token cache must still require an explicit token")
	}
	if _, err := manualClient.Get(context.Background(), "/resource", nil, larkcore.AccessTokenTypeUser,
		larkcore.WithUserAccessToken("test-user-token")); err != nil {
		t.Fatalf("explicit UAT with disabled cache: %v", err)
	}
	if atomic.LoadInt64(&tokenCalls) != 3 || atomic.LoadInt64(&provider.calls) != 3 || atomic.LoadInt64(&businessCalls) != 4 {
		t.Fatal("disabled cache must not automatically exchange credentials")
	}
}
