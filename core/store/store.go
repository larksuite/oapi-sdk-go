package store

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"sync"
	"time"
)

type Store interface {
	Get(context.Context, string) (string, error)
	Put(context.Context, string, string, time.Duration) error
}

func NewDefaultStoreWithLog(logger *log.LoggerProxy) *DefaultStore {
	return &DefaultStore{log: logger}
}

type DefaultStore struct {
	m   sync.Map
	log *log.LoggerProxy
}

type Value struct {
	value      string
	expireTime time.Time
}

func (s *DefaultStore) Get(ctx context.Context, key string) (string, error) {
	if val, ok := s.m.Load(key); ok {
		ev := val.(*Value)
		if ev.expireTime.After(time.Now()) {
			s.log.Debug(ctx, fmt.Sprintf("default store Get key %s, value is %s", key, ev.value))
			return ev.value, nil
		}
	}
	s.log.Debug(ctx, fmt.Sprintf("default store Get key %s, value is empty", key))
	return "", nil
}

func (s *DefaultStore) Put(ctx context.Context, key, value string, ttl time.Duration) error {
	expireTime := time.Now().Add(ttl)
	s.m.Store(key, &Value{
		value:      value,
		expireTime: expireTime,
	})
	s.log.Debug(ctx, fmt.Sprintf("default store put key %s, value is %s, expire time is %s", key, value, expireTime))
	return nil
}
