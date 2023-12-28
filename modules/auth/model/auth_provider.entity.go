package authmodel

import (
	"gogo/common"
)

type AuthProviderName string

const (
	GoogleAuthProvider AuthProviderName = "google"
)

type AuthProvider struct {
	common.SQLModel
	ProviderName AuthProviderName `json:"providerName" gorm:"column:provider_name;"`
	ProviderId   string           `json:"providerId" gorm:"column:provider_id;"`
	AuthId       string           `json:"authId" gorm:"column:auth_id;"`
}

func (AuthProvider) EntityName() string { return "auth_provider" }

func (AuthProvider) TableName() string { return "auth_provider" }
