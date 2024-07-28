package redisprovider

import (
	"context"
	"gogo/common"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	GetClient() *redis.Client
	SetValue(ctx context.Context, key string, value string, tls time.Duration) (string, error)
	DelValue(ctx context.Context, keys ...string) (int64, error)
	GetObjValue(ctx context.Context, key string, data interface{}) error
	GetStringValue(ctx context.Context, key string) (string, error)
}

type redisService struct {
	client *redis.Client
}

func NewRedisService(config redis.Options) RedisService {
	client := redis.NewClient(&config)
	return &redisService{client: client}
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

func (r *redisService) GetClient() *redis.Client {
	return r.client
}
