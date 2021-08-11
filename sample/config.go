package sample

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func TestConfigWithLogrusAndRedisStore(domain lark.Domain) core.Config {
	conf := lark.NewInternalAppConfigByEnv(domain)
	conf.SetLogger(Logrus{})
	conf.SetLogLevel(lark.LogLevelDebug)
	conf.SetStore(NewRedisStore())
	return conf
}

func TestConfig(domain lark.Domain) core.Config {
	conf := lark.NewInternalAppConfigByEnv(domain)
	conf.SetLogLevel(lark.LogLevelDebug)
	return conf
}

// use logrus implement lark.Logger

type Logrus struct {
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func (Logrus) Debug(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Info(ctx context.Context, args ...interface{}) {
	logrus.Info(args...)
}
func (Logrus) Warn(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Error(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}

// use redis implement lark.Store

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
