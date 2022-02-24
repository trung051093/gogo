package user

import (
	"user_management/common"
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
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(cond map[string]interface{}, userUpdate *usermodel.UserUpdate) error {
	if err := r.db.Where(cond).Updates(&userUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(cond map[string]interface{}) error {
	if err := r.db.Where(cond).Delete(&usermodel.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Get(id uint) (*usermodel.User, error) {
	var user *usermodel.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Search(cond map[string]interface{}, paging *common.Pagination) ([]usermodel.User, error) {
	var users []usermodel.User

	if err := r.db.Model(&usermodel.User{}).Where(cond).Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := r.db.Limit(paging.Limit).Offset(paging.Offset).Where(cond).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
