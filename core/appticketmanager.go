package core

import (
	"context"
	"fmt"
	"time"
)

var appTicketManager AppTicketManager = AppTicketManager{cache: cache}

func GetAppTicketManager() *AppTicketManager {
	return &appTicketManager
}

type AppTicketManager struct {
	cache Cache
}

func (m *AppTicketManager) Get(ctx context.Context, appId string) (string, error) {
	return m.cache.Get(ctx, appTicketKey(appId))
}

func (m *AppTicketManager) Put(ctx context.Context, appId, value string, ttl time.Duration) error {
	return m.cache.Set(ctx, appTicketKey(appId), value, ttl)
}

func appTicketKey(appID string) string {
	return fmt.Sprintf("%s-%s", appTicketKeyPrefix, appID)
}
