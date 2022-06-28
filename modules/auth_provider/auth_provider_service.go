package authprovider

import (
	"context"
	authprovidermodel "gogo/modules/auth_provider/model"
)

type authProviderService struct {
	repo AuthProviderRepository
}

func NewAuthProviderService(repo AuthProviderRepository) AuthProviderService {
	return &authProviderService{repo: repo}
}

func (s *authProviderService) Create(ctx context.Context, provider *authprovidermodel.AuthProviderCreate) (int, error) {
	return s.repo.Create(ctx, provider)
}

func (s *authProviderService) SearchOne(ctx context.Context, cond map[string]interface{}) (*authprovidermodel.AuthProvider, error) {
	return s.repo.SearchOne(ctx, cond)
}
