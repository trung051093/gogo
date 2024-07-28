package entity

import "github.com/dgrijalva/jwt-go"

type TokenPayload struct {
	Session  string `json:"session"`
	Username string `json:"username"`
	Role     string `json:"role"`
	AuthId   string `json:"auth_id"`
	UserId   string `json:"user_id"`
}

type TokenProvider struct {
	Session string `json:"session"`
	Token   string `json:"token"`
	Expiry  int64  `json:"expiry"`
	Created int64  `json:"created"`
}

type Claims struct {
	jwt.StandardClaims
	Payload TokenPayload `json:"payload"`
}
