package larkcore

import (
	"context"
	"testing"
	"time"
)

func TestTokenManagerSetAndGet(t *testing.T) {
	config := mockConfig()
	cache := &localCache{}
	tokenManager := TokenManager{cache: cache}

	err := tokenManager.set(context.Background(), tenantAccessTokenKey(config.AppId, "tenantKey"), "tokenValue", time.Minute)
	if err != nil {
		t.Errorf("set key failed ,%v", err)
	}

	token, err := tokenManager.getTenantAccessToken(context.Background(), config, "tenantKey")
	if err != nil {
		t.Errorf("get key failed ,%v", err)
	}

	if token == "" {
		t.Errorf("get key failed ,%v", err)
	}

}
