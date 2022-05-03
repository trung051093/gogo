package auth

import (
	"context"
	authmodel "user_management/modules/auth/model"
	authprovider "user_management/modules/auth_providers"
)

type AuthService interface {
	Login(ctx context.Context, payload *authmodel.AuthLogin) (*authprovider.TokenProvider, error)
	Register(ctx context.Context, payload authmodel.AuthRegister)
}
