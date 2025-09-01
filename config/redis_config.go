package config

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (c *RedisClient) Set(key, value string, expireTimeDuration time.Duration) error {
	return c.rdb.Set(c.ctx, key, value, expireTimeDuration).Err()
}

func (c *RedisClient) Get(key string) (string, error) {
	result, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) { //没有获取到值
			return "", nil
		} else {
			return "", err
		}
	}
	return result, nil
}

func (c *RedisClient) Close() error {
	return c.rdb.Close()
}

func (c *RedisClient) Ping() error {
	return c.rdb.Ping(c.ctx).Err()
}
