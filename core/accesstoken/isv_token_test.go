package accesstoken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/accesstoken/authorizationcode"
	"github.com/larksuite/oapi-sdk-go/v3/core/accesstoken/refreshtoken"
)

func TestISVUserTokenJSON(t *testing.T) {
	for _, grant := range []string{larkcore.GrantTypeAuthorizationCode, larkcore.GrantTypeRefreshToken} {
		t.Run(grant, func(t *testing.T) {
			var body map[string]string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
					t.Error("expected JSON OAuth body")
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode user token body: %v", err)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(accessTokenResponseBody{
					AccessToken: "test-user-token", TokenType: "Bearer", ExpiresIn: 7200,
					RefreshToken: "test-next-refresh-token",
				}); err != nil {
					t.Errorf("encode user token response: %v", err)
				}
			}))
			defer server.Close()
			provider := &mockClientAssertionProvider{token: &larkcore.Token{Value: "test+assertion&value"}}
			config := newTestConfig(server, provider)
			config.AppType = larkcore.AppTypeMarketplace
			config.AppSecret = "unused-secret"
			config.ReqTimeout = time.Second
			accessToken := NewAccessToken(config)
			var response *AccessTokenResp
			var err error
			if grant == larkcore.GrantTypeAuthorizationCode {
				response, err = accessToken.RetrieveByAuthorizationCode(context.Background(), authorizationcode.NewTokenRequestBuilder().
					Code("test+code&value").RedirectUri("https://example.com/callback?state=a+b&x=y").
					CodeVerifier("test-verifier").Scope("user_profile offline_access").Build(), larkcore.WithTenantKey("ignored-tenant"))
			} else {
				response, err = accessToken.Refresh(context.Background(), refreshtoken.NewTokenRequestBuilder().
					RefreshToken("test+refresh&value").Scope("user_profile").Build(), larkcore.WithTenantKey("ignored-tenant"))
			}
			if err != nil {
				t.Fatalf("exchange user token: %v", err)
			}
			if larkcore.StringValue(response.Data.AccessToken) != "test-user-token" ||
				larkcore.StringValue(response.Data.RefreshToken) != "test-next-refresh-token" {
				t.Fatal("unexpected user token response")
			}
			if body["grant_type"] != grant || body["client_id"] != config.AppId ||
				body["client_assertion_type"] != larkcore.ClientAssertionTypeJWTBearer ||
				body["client_assertion"] != provider.token.Value {
				t.Fatal("incorrect application authentication fields")
			}
			for _, field := range []string{"tenant_key", "client_secret", "app_secret", "app_ticket", "app_access_token"} {
				if _, ok := body[field]; ok {
					t.Errorf("unexpected field in UAT request: %s", field)
				}
			}
			if grant == larkcore.GrantTypeAuthorizationCode {
				if body["code"] != "test+code&value" || body["redirect_uri"] != "https://example.com/callback?state=a+b&x=y" ||
					body["code_verifier"] != "test-verifier" || body["scope"] != "user_profile offline_access" {
					t.Fatal("authorization parameters were not preserved")
				}
			} else if body["refresh_token"] != "test+refresh&value" || body["scope"] != "user_profile" {
				t.Fatal("refresh parameters were not preserved")
			}
		})
	}
}
