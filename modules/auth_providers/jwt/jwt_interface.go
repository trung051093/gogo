package jwtauthprovider

import authprovider "user_management/modules/auth_providers"

type JWTProvider interface {
	Generate(data authprovider.TokenPayload, expired uint) (*authprovider.TokenProvider, error)
	Validate(token string) (*authprovider.TokenPayload, error)
}
