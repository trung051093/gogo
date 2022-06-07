package redisprovider

import (
	"context"
	"time"
	"user_management/common"

	"github.com/go-redis/redis/v8"
)

type RedisService interface {
	SetValue(ctx context.Context, key string, value string, tls time.Duration) (string, error)
	DelValue(ctx context.Context, keys ...string) (int64, error)
	GetObjValue(ctx context.Context, key string, data interface{}) error
	GetStringValue(ctx context.Context, key string) (string, error)
}

type redisService struct {
	config *redis.Options
	client *redis.Client
}

func NewRedisService(config *redis.Options) *redisService {
	client := redis.NewClient(config)
	return &redisService{client: client, config: config}
}

func (r *redisService) SetValue(ctx context.Context, key string, value string, tls time.Duration) (string, error) {
	return r.client.Set(ctx, key, value, tls).Result()
}

func (r *redisService) DelValue(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

func (r *redisService) GetObjValue(ctx context.Context, key string, data interface{}) error {
	str, geterr := r.client.Get(ctx, key).Result()
	if geterr != nil {
		return geterr
	}

	parseErr := common.StringToJson(str, data)
	if parseErr != nil {
		return parseErr
	}

	return nil
}

func (r *redisService) GetStringValue(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
