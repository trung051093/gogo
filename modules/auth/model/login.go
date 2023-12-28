package authmodel

import jwtauthprovider "gogo/modules/auth/providers/jwt"

type AuthLoginDto struct {
	UserName string `json:"username" validate:"required"`
	Password string ` json:"password" validate:"required"`
}

type AuthResponseDto struct {
	jwtauthprovider.TokenProvider
}
