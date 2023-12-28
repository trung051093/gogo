package auth

import (
	"gogo/common"
	authmodel "gogo/modules/auth/model"

	"gorm.io/gorm"
)

type AuthProviderRepository interface {
	common.Repository[authmodel.AuthProvider]
}

type authProviderRepository struct {
	common.Repository[authmodel.AuthProvider]
}

func NewAuthProviderRepository(db *gorm.DB) AuthProviderRepository {
	repository := common.NewRepository[authmodel.AuthProvider](db)
	return &authProviderRepository{repository}
}
