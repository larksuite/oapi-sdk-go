package core

import (
	"context"
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	cache := localCache{}
	err := cache.Set(context.Background(), "key1", "value1", time.Second)
	if err != nil {
		t.Errorf("set key failed ,%v", err)
	}

	token, err := cache.Get(context.Background(), "key1")
	if err != nil {
		t.Errorf("get key failed ,%v", err)
	}
	if token == "" {
		t.Errorf("get key empty ,%v", err)

	}
}

//
//func TestCacheTimeout(t *testing.T) {
//
//	cache := localCache{}
//	err := cache.Set(context.Background(), "key1", "value1", time.Second)
//	if err != nil {
//		t.Errorf("set key failed ,%v", err)
//	}
//
//	time.Sleep(2 * time.Second)
//
//	token, err := cache.Get(context.Background(), "key1")
//	if err != nil {
//		t.Errorf("get key failed ,%v", err)
//	}
//	if token == "" {
//		t.Errorf("get key empty ,%v", err)
//
//	}
//}
