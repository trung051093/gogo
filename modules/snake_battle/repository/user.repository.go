package repository

import (
	"gogo/common"
	"gogo/modules/snake_battle/entity"

	"gorm.io/gorm"
)

// User corresponds to the users table
type User struct {
	entity.User
}

func (User) TableName() string {
	return entity.User{}.TableName()
}

type UserRepository interface {
	common.Repository[User]
	// Add other specific user methods if needed, e.g., FindByUsername
}

type userRepository struct {
	common.Repository[User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	repo := common.NewRepository[User](db)
	// Set preloads if needed
	return &userRepository{repo}
}
