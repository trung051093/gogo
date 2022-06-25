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

// expireTime: 30d, 24h,...
func (s *authService) setToken(ctx context.Context, key string, token string, expireDays int) error {
	ttl := DefaultTTL
	if paramsTTL, err := time.ParseDuration(fmt.Sprintf("%dd", expireDays)); err != nil {
		ttl = paramsTTL
	}
	return s.cacheService.Set(&cache.Item{
		Key:   key,
		TTL:   ttl,
		Value: new(string),
		Do: func(i *cache.Item) (interface{}, error) {
			return token, nil
		},
	})
}

func (s *authService) getToken(ctx context.Context, key string) (string, error) {
	var token string
	err := s.cacheService.Get(ctx, key, &token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) deleteToken(ctx context.Context, key string) error {
	return s.cacheService.Delete(ctx, key)
}

func (s *authService) getKeyToken(ctx context.Context, keys ...string) string {
	return strings.Join(keys, ":")
}
