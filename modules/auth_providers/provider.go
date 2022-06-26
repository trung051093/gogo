package authprovider

import (
	"errors"
	"gogo/common"
)

var (
	ErrorInvalidToken     = common.NewCustomError(errors.New("Invalid token"), "Invalid token", "ErrorInvalidToken")
	ErrorInvalidSignature = common.NewCustomError(errors.New("Invalid signature"), "Invalid signature", "ErrorInvalidSignature")
)

type TokenPayload struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type TokenProvider struct {
	Token   string `json:"token"`
	Expiry  int64  `json:"expiry"`
	Created int64  `json:"created"`
}
