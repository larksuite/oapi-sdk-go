package lark

import (
	"context"
	"sync"
	"time"
)

type store interface {
	Get(context.Context, string) (string, error)
	Put(context.Context, string, string, time.Duration) error
}

type defaultStore struct {
	m sync.Map
}

type Value struct {
	value      string
	expireTime time.Time
}

func (s *defaultStore) Get(ctx context.Context, key string) (string, error) {
	if val, ok := s.m.Load(key); ok {
		ev := val.(*Value)
		if ev.expireTime.After(time.Now()) {
			return ev.value, nil
		}
	}
	return "", nil
}

func (s *defaultStore) Put(ctx context.Context, key, value string, ttl time.Duration) error {
	expireTime := time.Now().Add(ttl)
	s.m.Store(key, &Value{
		value:      value,
		expireTime: expireTime,
	})
	return nil
}
