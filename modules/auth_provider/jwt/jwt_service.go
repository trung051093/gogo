package jwtauthprovider

import (
	"gogo/common"
	authmodelprovider "gogo/modules/auth_provider/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	Payload authmodelprovider.TokenPayload `json:"payload"`
}

type jwtProvider struct {
	secret string
}

func NewJWTProvider(secret string) *jwtProvider {
	return &jwtProvider{secret: secret}
}

func (jwtP *jwtProvider) Generate(data authmodelprovider.TokenPayload, expired uint) (*authmodelprovider.TokenProvider, error) {
	expiresAt := time.Now().Add(time.Duration(expired) * time.Hour * 24).Unix()
	createdAt := time.Now().Unix()
	claims := &Claims{
		Payload: data,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtP.secret))
	if err != nil {
		return nil, err
	}
	return &authmodelprovider.TokenProvider{
		Token:   tokenString,
		Expiry:  expiresAt,
		Created: createdAt,
	}, nil
}

func (jwtP *jwtProvider) Validate(token string) (*authmodelprovider.TokenPayload, error) {
	claims := &Claims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtP.secret), nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, common.ErrorUnauthorized()
	}
	return &claims.Payload, nil
}
