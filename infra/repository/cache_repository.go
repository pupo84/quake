package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

type CacheRepositorier interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type CacheRepository struct {
	logger *zap.SugaredLogger
	redis  *redis.Client
}

func NewCacheRepository(logger *zap.SugaredLogger) *CacheRepository {
	redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", viper.GetString("REDIS_ADDRESS"), viper.GetInt("REDIS_PORT")),
		DB:   viper.GetInt("REDIS_DB"),
	})
	return &CacheRepository{logger: logger, redis: redis}
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}

func (r *CacheRepository) Set(ctx context.Context, key string, value string) error {
	return r.redis.Set(ctx, key, value, 0).Err()
}
