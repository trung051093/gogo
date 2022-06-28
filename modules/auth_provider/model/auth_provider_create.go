package authprovidermodel

import "gogo/common"

type AuthProviderCreate struct {
	common.SQLModel
	ProviderName string `json:"providerName" gorm:"column:provider_name;"`
	ProviderId   string `json:"providerId" gorm:"column:provider_id;"`
	UserId       int    `json:"userId" gorm:"column:user_id;"`
}

func (AuthProviderCreate) TableName() string { return AuthProvider{}.TableName() }
