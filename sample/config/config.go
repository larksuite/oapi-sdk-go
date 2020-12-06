package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/sirupsen/logrus"
	"time"
)

func GetConfig(domain constants.Domain, appSettings *config.AppSettings, level log.Level) *config.Config {
	logger := Logrus{}
	store := NewRedisStore()
	coreCtx := core.WrapContext(context.Background())
	coreCtx.GetHTTPStatusCode()
	return config.NewConfig(constants.DomainLarkSuite, appSettings, logger, level, store)
}

// use logrus implement log.Logger
type Logrus struct {
}

func (Logrus) Debug(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Info(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Warn(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Error(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}

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
