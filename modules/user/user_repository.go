package user

import (
	"time"
	usermodel "user_management/modules/user/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *usermodel.UserCreate) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	r.db.Model(user).Update("CreatedAt", time.Now())
	return nil
}

func (r *UserRepository) Update(cond map[string]interface{}, userUpdate *usermodel.UserUpdate) error {
	if err := r.db.Where(cond).Updates(userUpdate).Error; err != nil {
		return err
	}
	r.db.Model(userUpdate).Update("UpdatedAt", time.Now())
	return nil
}

func (r *UserRepository) Delete(cond map[string]interface{}) (*usermodel.User, error) {
	var user *usermodel.User
	if err := r.db.Where(cond).Updates("status = 0").Error; err != nil {
		return nil, err
	}
	r.db.Model(&user).Update("DeletedAt", time.Now())
	return user, nil
}

func (r *UserRepository) Get(id uint) (*usermodel.User, error) {
	var user *usermodel.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Search(cond map[string]interface{}) ([]usermodel.User, error) {
	var users []usermodel.User

	if err := r.db.Where(cond).Find(users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
