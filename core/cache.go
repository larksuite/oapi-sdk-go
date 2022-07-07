package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/larksuite/oapi-sdk-go/httpclient"
)

var cache = &localCache{}

func NewCache(config *Config) {
	if config.TokenCache != nil {
		tokenManager = TokenManager{cache: config.TokenCache}
		appTicketManager = AppTicketManager{cache: config.TokenCache}
	}
}

func NewHttpClient(config *Config) {
	if config.HttpClient == nil {
		config.HttpClient = httpclient.NewHttpClient(config.ReqTimeout)
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
		//TODO del
		fmt.Println(fmt.Sprintf("get key:%s,hit cache,time left %f seconds",
			key, ev.expireTime.Sub(time.Now()).Seconds()))
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
