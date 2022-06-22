package cacheprovider

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	Once(item *cache.Item) error
	Get(ctx context.Context, key string, value interface{}) error
	Set(item *cache.Item) error
}

type CacheConfig struct {
	Addrs map[string]string
}

type cacheService struct {
	ring    *redis.Ring
	mycache *cache.Cache
}

func NewCacheService(config *CacheConfig) *cacheService {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: config.Addrs,
	})

	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &cacheService{ring: ring, mycache: mycache}
}

func (s *cacheService) Once(item *cache.Item) error {
	return s.mycache.Once(item)
}

func (s *cacheService) Get(ctx context.Context, key string, value interface{}) error {
	return s.mycache.Get(ctx, key, value)
}

func (s *cacheService) Set(item *cache.Item) error {
	return s.mycache.Set(item)
}
