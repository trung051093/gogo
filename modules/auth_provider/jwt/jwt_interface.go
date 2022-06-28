package jwtauthprovider

import authmodelprovider "gogo/modules/auth_provider/model"

type JWTProvider interface {
	Generate(data authmodelprovider.TokenPayload, expired uint) (*authmodelprovider.TokenProvider, error)
	Validate(token string) (*authmodelprovider.TokenPayload, error)
}
