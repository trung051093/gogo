package user

import (
	"context"
	"gogo/common"
	usermodel "gogo/modules/user/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	common.Repository[usermodel.User]
	Deactive(ctx context.Context, user *usermodel.User) (*usermodel.User, error)
}

type userRepository struct {
	common.Repository[usermodel.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	repository := common.NewRepository[usermodel.User](db)
	return &userRepository{repository}
}

func (r *userRepository) Deactive(ctx context.Context, user *usermodel.User) (*usermodel.User, error) {
	return r.Updates(ctx, user, map[string]interface{}{"is_active": false})
}
