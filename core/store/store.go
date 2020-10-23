package store

import (
	"context"
	"log"
	"sync"
	"time"
)

type Store interface {
	Get(context.Context, string) (string, error)
	Put(context.Context, string, string, time.Duration) error
}

func NewDefaultStore() *DefaultStore {
	return &DefaultStore{}
}

type DefaultStore struct {
	m sync.Map
}

type Value struct {
	value      string
	expireTime time.Time
}

func (s *DefaultStore) Get(ctx context.Context, key string) (string, error) {
	if val, ok := s.m.Load(key); ok {
		ev := val.(*Value)
		if ev.expireTime.After(time.Now()) {
			log.Printf("default store Get key(%s), value is (%s)", key, ev.value)
			return ev.value, nil
		}
	}
	log.Printf("default store Get key(%s), value is empty", key)
	return "", nil
}

func (s *DefaultStore) Put(ctx context.Context, key, value string, ttl time.Duration) error {
	s.m.Store(key, &Value{
		value:      value,
		expireTime: time.Now().Add(ttl),
	})
	log.Printf("default store put key(%s), value is (%s)", key, value)
	return nil
}
