package authprovidermodel

import (
	"errors"
	"gogo/common"
)

const EntityName = "auth_provider"

const (
	GoogleAuthProvider = "google"
)

var (
	ErrorInvalidToken     = common.NewCustomError(errors.New("invalid token"), "Invalid token", "ErrorInvalidToken")
	ErrorInvalidSignature = common.NewCustomError(errors.New("invalid signature"), "Invalid signature", "ErrorInvalidSignature")
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

type AuthProvider struct {
	common.SQLModel
	ProviderName string `json:"providerName" gorm:"column:provider_name;"`
	ProviderId   string `json:"providerId" gorm:"column:provider_id;"`
	UserId       int    `json:"userId" gorm:"column:user_id;"`
}

func (AuthProvider) TableName() string { return "auth_provider" }
