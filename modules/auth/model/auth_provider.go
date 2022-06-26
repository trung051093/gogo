package authmodel

import (
	"gogo/common"
)

const EntityName = "auth"

type AuthProvider struct {
	common.SQLModel
	ProviderName string `json:"providerName" gorm:"column:provider_name;"`
	ProviderId   string `json:"providerId" gorm:"column:provider_id;"`
}
