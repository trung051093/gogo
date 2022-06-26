package jwtauthprovider

import authprovider "gogo/modules/auth_providers"

type JWTProvider interface {
	Generate(data authprovider.TokenPayload, expired uint) (*authprovider.TokenProvider, error)
	Validate(token string) (*authprovider.TokenPayload, error)
}
