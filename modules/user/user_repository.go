package user

import (
	"context"
	"fmt"
	"user_management/common"
	usermodel "user_management/modules/user/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *usermodel.UserCreate) (int, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return -1, common.ErrorCannotCreateEntity(usermodel.EntityName, err)
	}
	return user.Id, nil
}

func (r *userRepository) Update(ctx context.Context, id int, userUpdate *usermodel.UserUpdate) (int, error) {
	if err := r.db.WithContext(ctx).Where(map[string]interface{}{"id": id}).Updates(&userUpdate).Error; err != nil {
		return -1, common.ErrorCannotUpdateEntity(usermodel.EntityName, err)
	}
	return userUpdate.Id, nil
}

func (r *userRepository) DeActive(ctx context.Context, user *usermodel.User) (int, error) {
	if err := r.db.WithContext(ctx).Where(map[string]interface{}{"id": user.Id}).Updates(map[string]interface{}{"is_active": 0}).Error; err != nil {
		return -1, common.ErrorCannotDeleteEntity(usermodel.EntityName, err)
	}
	return user.Id, nil
}

func (r *userRepository) Delete(ctx context.Context, user *usermodel.User) (int, error) {
	if err := r.db.WithContext(ctx).Delete(&user).Error; err != nil {
		return -1, common.ErrorCannotDeleteEntity(usermodel.EntityName, err)
	}
	return user.Id, nil
}

func (r *userRepository) Get(ctx context.Context, id int) (*usermodel.User, error) {
	var user *usermodel.User

	if err := r.db.WithContext(ctx).Model(&usermodel.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, common.ErrorCannotFoundEntity(usermodel.EntityName, err)
	}

	return user, nil
}

func (r *userRepository) Search(ctx context.Context, cond map[string]interface{}, filter *usermodel.UserFilter, paging *common.Pagination) ([]usermodel.User, error) {
	var users []usermodel.User
	order := fmt.Sprintf("%s %s", filter.SortField, filter.SortName)
	if err := r.db.WithContext(ctx).Model(&usermodel.User{}).Where(cond).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrorCannotListEntity(usermodel.EntityName, err)
	}

	if err := r.db.WithContext(ctx).Model(&usermodel.User{}).Limit(paging.Limit).Offset(paging.Offset).Where(cond).Order(order).Select(filter.Fields).Find(&users).Error; err != nil {
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
