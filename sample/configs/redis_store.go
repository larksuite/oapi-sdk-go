package configs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// use redis implement store.Store
type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &RedisStore{client: rdb}
}

func (rs *RedisStore) Get(ctx context.Context, key string) (string, error) {
	if rs.client == nil {
		panic("redis client is not initialized")
	}
	value, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (rs *RedisStore) Put(ctx context.Context, key string, value string, expired time.Duration) error {
	if rs.client == nil {
		panic("redis client is not initialized")
	}
	_, err := rs.client.Set(ctx, key, value, expired).Result()
	if err != nil {
		return err
	}
	return nil
}
