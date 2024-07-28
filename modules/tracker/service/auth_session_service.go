package service

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/cache/v9"
)

const (
	DefaultTTL          = 30 * 24 * time.Hour
	AuthenticationToken = "AuthenticationToken"
	ForgotPasswordToken = "ForgotPasswordToken"
)

func (s *authService) deleteKey(ctx context.Context, key string) error {
	return s.cacheService.Delete(ctx, key)
}

func (s *authService) getCacheKey(keys ...string) string {
	return strings.Join(keys, ":")
}

func (s *authService) getAuthenicationKey(session string) string {
	return s.getCacheKey(AuthenticationToken, session)
}

// expireTime: 30d, 24h,...
func (s *authService) setJwtSession(_ context.Context, key string, jwt string) error {
	return s.cacheService.Set(&cache.Item{
		Key:   key,
		TTL:   DefaultTTL,
		Value: new(string),
		Do: func(i *cache.Item) (interface{}, error) {
			return jwt, nil
		},
	})
}

// expireTime: 30d, 24h,...
func (s *authService) getJwtSession(ctx context.Context, key string) (string, error) {
	var jwt string
	err := s.cacheService.Get(ctx, key, &jwt)
	if err != nil {
		return "", err
	}
	return jwt, nil
}
