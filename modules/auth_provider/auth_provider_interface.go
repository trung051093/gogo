package authprovider

import (
	"context"
	authprovidermodel "gogo/modules/auth_provider/model"
)

type Read interface {
	SearchOne(ctx context.Context, cond map[string]interface{}) (*authprovidermodel.AuthProvider, error)
}

type Write interface {
	Create(ctx context.Context, provider *authprovidermodel.AuthProviderCreate) (int, error)
}

type AuthProviderRepository interface {
	Read
	Write
}

type AuthProviderService interface {
	Create(ctx context.Context, provider *authprovidermodel.AuthProviderCreate) (int, error)
	SearchOne(ctx context.Context, cond map[string]interface{}) (*authprovidermodel.AuthProvider, error)
}
