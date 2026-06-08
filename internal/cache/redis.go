package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
	})
}

type redisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) Cache {
	return &redisStore{
		rdb: rdb,
	}
}

func (r *redisStore) Save(ctx context.Context, key, value string, duration time.Duration) error {
	return r.rdb.Set(ctx, key, value, duration).Err()
}

func (r *redisStore) Get(ctx context.Context, key string) (string, error) {
	value, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("cache miss")
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *redisStore) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}
