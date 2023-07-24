package repository

import (
	"context"
	"errors"
	"go-ekyc/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisConfig config.RedisConfig
	redisClient *redis.Client
}
func newRedisRepository(redisConfig config.RedisConfig) (*RedisRepository, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Endpoint,
		Password: "", 
		DB:       0,  
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New("error while connecting to redis")
	}
	return &RedisRepository{
		redisConfig: redisConfig,
		redisClient: redisClient,
	}, nil
}
func (r *RedisRepository) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ok, err := r.redisClient.SetNX(ctx, key, value, expiration).Result()

	return ok, err
}


