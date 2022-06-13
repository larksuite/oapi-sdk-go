package core

import (
	"context"
	"sync"
	"time"
)

var cache = &localCache{}

func NewCache(config *Config) {
	if config.TokenCache != nil {
		tokenManager = TokenManager{cache: config.TokenCache}
		appTicketManager = AppTicketManager{cache: config.TokenCache}
	}
}

type Cache interface {
	Set(ctx context.Context, key string, value string, expireTime time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type localCache struct {
	m sync.Map
}

func (s *localCache) Get(ctx context.Context, key string) (string, error) {
	if val, ok := s.m.Load(key); ok {
		ev := val.(*Value)
		if ev.expireTime.After(time.Now()) {
			return ev.value, nil
		}
	}
	return "", nil
}

func (s *localCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	expireTime := time.Now().Add(ttl)
	s.m.Store(key, &Value{
		value:      value,
		expireTime: expireTime,
	})
	return nil
}

type Value struct {
	value      string
	expireTime time.Time
}
