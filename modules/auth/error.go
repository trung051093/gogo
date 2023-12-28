package auth

import (
	"errors"
	"gogo/common"
)

var (
	ErrorInvalidToken     = common.NewCustomError(errors.New("invalid token"), "Invalid token", "ErrorInvalidToken")
	ErrorInvalidSignature = common.NewCustomError(errors.New("invalid signature"), "Invalid signature", "ErrorInvalidSignature")
)
