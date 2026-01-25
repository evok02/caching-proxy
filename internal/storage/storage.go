package storage

import (
	"context"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"os"
	"time"
)

type RedisStorage interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

func InitDB() (*redis.Client, error) {
	dbPort := os.Getenv("DB_PORT")
	fmt.Println(dbPort)
	rdb := redis.NewClient(&redis.Options{
		Addr: "cache:" + dbPort,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
