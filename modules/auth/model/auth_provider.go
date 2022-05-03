package authmodel

import (
	"user_management/common"
)

const EntityName = "auth"

type AuthProvider struct {
	common.SQLModel
	ProviderName string `json:"providerName,omitempty" gorm:"column:provider_name;"`
	ProviderId   string `json:"providerId,omitempty" gorm:"column:provider_id;"`
}
