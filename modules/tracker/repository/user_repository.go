package repository

import (
	"gogo/common"
	"gogo/modules/tracker/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	common.Repository[entity.User]
}

type userRepository struct {
	common.Repository[entity.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	repository := common.NewRepository[entity.User](db)
	return &userRepository{repository}
}
