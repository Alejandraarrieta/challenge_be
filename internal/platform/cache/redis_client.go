package cache

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClientInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	Close() error
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
}

var (
	redisClient RedisClientInterface
	once        sync.Once
)

func initRedisClient() error {
	addr := os.Getenv("REDIS_ADDR")     // Ejemplo: "localhost:6379"
	password := os.Getenv("REDIS_PASS") // Puede estar vacío
	db := 0                             // Por defecto usamos DB 0

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

func GetRedisClientInstance() (RedisClientInterface, error) {
	var err error
	once.Do(func() {
		err = initRedisClient()
	})

	if err != nil {
		return nil, err
	}

	return redisClient, nil
}

func Close() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
