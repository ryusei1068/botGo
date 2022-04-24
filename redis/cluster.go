package redis

import (
	"context"

	"github.com/botGo/config"
	"github.com/go-redis/redis/v8"
)

type (
	RedisClient struct {
		redisClient *redis.Client
		ctx         context.Context
	}

	IRedis interface {
		SetValue(string, string) error
		GetValues(string) []string
		GetKeys(string) []string
	}
)

func (rc *RedisClient) SetValue(key string, value string) error {
	err := rc.redisClient.Set(rc.ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) GetValues(key string) []string {
	return []string{}
}

func (rc *RedisClient) GetKeys(key string) []string {
	return []string{}
}

func NewRedisClient() IRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GoDotEnvVariable("ADDR"),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	return &RedisClient{redisClient: rdb, ctx: ctx}
}
