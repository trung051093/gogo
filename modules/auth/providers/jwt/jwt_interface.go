package jwtauthprovider

type TokenPayload struct {
	Email  string `json:"email"`
	Role   string `json:"role"`
	AuthId string `json:"authId"`
	UserId string `json:"userId"`
}

type TokenProvider struct {
	Token   string `json:"token"`
	Expiry  int64  `json:"expiry"`
	Created int64  `json:"created"`
}

type JWTProvider interface {
	Generate(data TokenPayload, expired uint) (*TokenProvider, error)
	Validate(token string) (*TokenPayload, error)
}
