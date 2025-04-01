package repository

import (
	"gogo/common"
	"gogo/modules/snake_battle/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	common.Repository[entity.User]
	// Add other specific user methods if needed, e.g., FindByUsername
}

type userRepository struct {
	common.Repository[entity.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	repo := common.NewRepository[entity.User](db)
	// Set preloads if needed
	return &userRepository{repo}
}
