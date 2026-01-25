package storage

import (
	"context"
	"github.com/evok02/cacher/internal/config"
	redis "github.com/redis/go-redis/v9"
	"time"
)

type RedisStorage interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

func InitDB(cfg config.DBConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Name + ":" + cfg.Port,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
