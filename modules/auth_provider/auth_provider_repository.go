package authprovider

import (
	"context"
	"gogo/common"
	authprovidermodel "gogo/modules/auth_provider/model"

	"gorm.io/gorm"
)

type authProviderRepository struct {
	db *gorm.DB
}

func NewAuthProviderRepository(db *gorm.DB) AuthProviderRepository {
	return &authProviderRepository{db: db}
}

func (r *authProviderRepository) Create(ctx context.Context, provider *authprovidermodel.AuthProviderCreate) (int, error) {
	if err := r.db.WithContext(ctx).Create(&provider).Error; err != nil {
		return -1, common.ErrorCannotCreateEntity(authprovidermodel.EntityName, err)
	}
	return provider.Id, nil
}

func (r *authProviderRepository) SearchOne(ctx context.Context, cond map[string]interface{}) (*authprovidermodel.AuthProvider, error) {
	var provider *authprovidermodel.AuthProvider

	if err := r.db.WithContext(ctx).Model(&authprovidermodel.AuthProvider{}).Where(cond).First(&provider).Error; err != nil {
		return nil, common.ErrorCannotCreateEntity(authprovidermodel.EntityName, err)
	}
	return provider, nil
}
