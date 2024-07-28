package cacheprovider

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Once(item *cache.Item) error
	Get(ctx context.Context, key string, value interface{}) error
	Set(item *cache.Item) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
}

type CacheConfig struct {
	Addrs map[string]string
}

type cacheService struct {
	client  *redis.Client
	mycache *cache.Cache
}

func NewCacheService(client *redis.Client) CacheService {
	mycache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &cacheService{client: client, mycache: mycache}
}

func (s *cacheService) Once(item *cache.Item) error {
	return s.mycache.Once(item)
}

func (s *cacheService) Set(item *cache.Item) error {
	return s.mycache.Set(item)
}

func (s *cacheService) Get(ctx context.Context, key string, value interface{}) error {
	return s.mycache.Get(ctx, key, value)
}

func (s *cacheService) Exists(ctx context.Context, key string) bool {
	return s.mycache.Exists(ctx, key)
}

func (s *cacheService) Delete(ctx context.Context, key string) error {
	return s.mycache.Delete(ctx, key)
}
