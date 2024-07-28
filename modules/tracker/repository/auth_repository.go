package repository

import (
	"context"
	"gogo/common"
	"gogo/modules/tracker/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	common.Repository[entity.Auth]
	FindByUserId(ctx context.Context, userId string) (*entity.Auth, error)
	FindByUserName(ctx context.Context, username string) (*entity.Auth, error)
}

type authRepository struct {
	common.Repository[entity.Auth]
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	repository := common.NewRepository[entity.Auth](db)
	repository.SetPreloadKeys("User")
	return &authRepository{repository}
}

func (r *authRepository) FindByUserId(ctx context.Context, userId string) (*entity.Auth, error) {
	query := common.NewGormQuery()
	query.SetQuery("user_id = ?", userId)
	return r.FindOneByCond(ctx, query)
}

func (r *authRepository) FindByUserName(ctx context.Context, username string) (*entity.Auth, error) {
	query := common.NewGormQuery()
	query.SetQuery("username = ?", username)
	return r.FindOneByCond(ctx, query)
}
