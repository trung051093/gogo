package user

import (
	"context"
	"user_management/common"
	usermodel "user_management/modules/user/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *usermodel.UserCreate) (int, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return -1, common.ErrorCannotCreateEntity(usermodel.EntityName, err)
	}
	return user.Id, nil
}

func (r *userRepository) Update(ctx context.Context, cond map[string]interface{}, userUpdate *usermodel.UserUpdate) error {
	if err := r.db.WithContext(ctx).Where(cond).Updates(&userUpdate).Error; err != nil {
		return common.ErrorCannotUpdateEntity(usermodel.EntityName, err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, cond map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Where(cond).Delete(&usermodel.User{}).Error; err != nil {
		return common.ErrorCannotDeleteEntity(usermodel.EntityName, err)
	}
	return nil
}

func (r *userRepository) Get(ctx context.Context, id uint) (*usermodel.User, error) {
	var user *usermodel.User

	if err := r.db.WithContext(ctx).Model(&usermodel.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, common.ErrorCannotFoundEntity(usermodel.EntityName, err)
	}

	return user, nil
}

func (r *userRepository) Search(ctx context.Context, cond map[string]interface{}, filter *usermodel.UserFilter, paging *common.Pagination) ([]usermodel.User, error) {
	var users []usermodel.User

	if err := r.db.WithContext(ctx).Model(&usermodel.User{}).Where(cond).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrorCannotListEntity(usermodel.EntityName, err)
	}

	if err := r.db.WithContext(ctx).Model(&usermodel.User{}).Limit(paging.Limit).Offset(paging.Offset).Where(cond).Order(filter.Order).Select(filter.Fields).Find(&users).Error; err != nil {
		return nil, common.ErrorCannotListEntity(usermodel.EntityName, err)
	}

	return users, nil
}

func (r *userRepository) SearchOne(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error) {
	var user *usermodel.User

	if err := r.db.WithContext(ctx).Model(&usermodel.User{}).Where(cond).First(&user).Error; err != nil {
		return nil, common.ErrorCannotFoundEntity(usermodel.EntityName, err)
	}

	return user, nil
}
