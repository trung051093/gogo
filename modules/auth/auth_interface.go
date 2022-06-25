package auth

import (
	"context"
	authmodel "user_management/modules/auth/model"
	authprovider "user_management/modules/auth_providers"
	usermodel "user_management/modules/user/model"
)

type AuthService interface {
	Register(ctx context.Context, payload *authmodel.AuthRegister) (int, error)
	Login(ctx context.Context, payload *authmodel.AuthLogin) (*authprovider.TokenProvider, error)
	Logout(ctx context.Context, user *usermodel.User) (int, error)
	ForgotPassword(ctx context.Context, payload *authmodel.AuthForgotPassword) (int, error)
	ResetPassword(ctx context.Context, payload *authmodel.AuthResetPassword) (int, error)
}
