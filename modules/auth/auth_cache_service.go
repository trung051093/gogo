package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/cache/v8"
)

const (
	DefaultTTL          = 24 * time.Hour
	AuthenticationToken = "AuthenticationToken"
	ForgotPasswordToken = "ForgotPasswordToken"
)

func (s *authService) deleteKey(ctx context.Context, key string) error {
	return s.cacheService.Delete(ctx, key)
}

func (s *authService) getCacheKey(ctx context.Context, keys ...string) string {
	return strings.Join(keys, ":")
}

// expireTime: 30d, 24h,...
func (s *authService) setSession(ctx context.Context, key string, data interface{}, expireDays int) error {
	ttl := DefaultTTL
	if paramsTTL, err := time.ParseDuration(fmt.Sprintf("%dd", expireDays)); err != nil {
		ttl = paramsTTL
	}
	return s.cacheService.Set(&cache.Item{
		Key:   key,
		TTL:   ttl,
		Value: new(interface{}),
		Do: func(i *cache.Item) (interface{}, error) {
			return data, nil
		},
	})
}

// expireTime: 30d, 24h,...
func (s *authService) getSession(ctx context.Context, key string) (interface{}, error) {
	var data interface{}
	err := s.cacheService.Get(ctx, key, &data)
	if err != nil {
		return "", err
	}
	return data, nil
}
