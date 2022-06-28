package auth

import (
	"context"
	authmodel "gogo/modules/auth/model"
	authprovidermodel "gogo/modules/auth_provider/model"
	usermodel "gogo/modules/user/model"
)

type AuthService interface {
	Register(ctx context.Context, payload *authmodel.AuthRegister) (int, error)
	Login(ctx context.Context, payload *authmodel.AuthLogin) (*authprovidermodel.TokenProvider, error)
	Logout(ctx context.Context, user *usermodel.User) (int, error)
	ForgotPassword(ctx context.Context, payload *authmodel.AuthForgotPassword) (int, error)
	ResetPassword(ctx context.Context, payload *authmodel.AuthResetPassword) (int, error)

	GoogleLogin(ctx context.Context, state string) string
	GoogleValidate(ctx context.Context, code string) (int, error)
}
