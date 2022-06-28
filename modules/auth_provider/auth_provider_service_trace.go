package authprovider

import (
	"context"
	decorator "gogo/decorators"
	authprovidermodel "gogo/modules/auth_provider/model"
)

func (s *authProviderService) SearchOneTrace(ctx context.Context, cond map[string]interface{}) (*authprovidermodel.AuthProvider, error) {
	data, err := decorator.TraceService[*authprovidermodel.AuthProvider](ctx, "authProviderService.SearchOne")(s, "SearchOne")(ctx, cond)
	return data, err
}

func (s *authProviderService) CreateTrace(ctx context.Context, provider *authprovidermodel.AuthProviderCreate) (int, error) {
	data, err := decorator.TraceService[int](ctx, "authProviderService.Create")(s, "Create")(ctx, provider)
	return data, err
}
