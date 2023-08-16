package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// CacheRepositorier is the interface that wraps the basic Get and Set methods
type CacheRepositorier interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

// CacheRepository is the struct that implements the CacheRepositorier interface
type CacheRepository struct {
	logger *zap.SugaredLogger
	redis  *redis.Client
}

// NewCacheRepository returns a new CacheRepository
func NewCacheRepository(logger *zap.SugaredLogger) *CacheRepository {
	redis := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", viper.GetString("REDIS_SERVER"), viper.GetInt("REDIS_PORT")),
		DB:   viper.GetInt("REDIS_DB"),
	})

	if _, err := redis.Ping(context.Background()).Result(); err != nil {
		logger.Fatalf("Error connecting to Redis: %s", err)
	}

	return &CacheRepository{logger: logger, redis: redis}
}

// Get returns the value of a given key
func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}

// Set sets the value of a given key
func (r *CacheRepository) Set(ctx context.Context, key string, value string) error {
	return r.redis.Set(ctx, key, value, 0).Err()
}
