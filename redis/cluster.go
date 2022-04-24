package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type (
	RedisClient struct {
		redisClient *redis.Client
		ctx         context.Context
	}

	IRedis interface {
		SetValues(string, []DirectInfo) error
		GetValues(string) ([]DirectInfo, error)
		GetKeys(string) []string
		getRange(string) int
	}

	DirectInfo struct {
		WebhookUrl       string `json:"webhookurl"`
		WebhookId        string `json:"webhookId"`
		WebhookToken     string `json:"webhooktoken"`
		DiscordChannelId string `json:"discordChannleid"`
	}
)

func (rc *RedisClient) SetValues(key string, values []DirectInfo) error {
	for i := range values {
		json, err := json.Marshal(&values[i])
		if err != nil {
			log.Printf("failed to convert type ([]byte), %v", err)
			return err
		}

		err = rc.redisClient.RPush(rc.ctx, key, json).Err()
		if err != nil {
			log.Printf("failed to add data, %v", err)
			return err
		}
	}

	return nil
}

func (rc *RedisClient) getRange(key string) int {
	ldata := rc.redisClient.LLen(rc.ctx, key)
	return int(ldata.Val())
}

func (rc *RedisClient) GetValues(key string) ([]DirectInfo, error) {
	endIndex := rc.getRange(key)
	val, err := rc.redisClient.LRange(rc.ctx, key, 0, int64(endIndex)).Result()

	if err != nil {
		log.Printf("failed to get values, %v", err)
		return nil, err
	}

	var directInfo DirectInfo
	var directInfoArr []DirectInfo

	for i := range val {
		err = json.Unmarshal([]byte(val[i]), &directInfo)

		if err != nil {
			log.Printf("failed to convert to json, %v", err)
			return nil, err
		}
		directInfoArr = append(directInfoArr, directInfo)
	}

	return directInfoArr, nil
}

func (rc *RedisClient) GetKeys(key string) []string {
	return []string{}
}

func NewRedisClient(addr string) IRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	return &RedisClient{redisClient: rdb, ctx: ctx}
}
