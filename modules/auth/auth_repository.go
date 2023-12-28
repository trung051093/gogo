package auth

import (
	"context"
	"gogo/common"
	authmodel "gogo/modules/auth/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	common.Repository[authmodel.Auth]
	FindByUserId(ctx context.Context, userId string) (*authmodel.Auth, error)
	FindByEmail(ctx context.Context, email string) (*authmodel.Auth, error)
	FindByPhone(ctx context.Context, phone string) (*authmodel.Auth, error)
	FindByUserName(ctx context.Context, username string) (*authmodel.Auth, error)
}

type authRepository struct {
	common.Repository[authmodel.Auth]
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	repository := common.NewRepository[authmodel.Auth](db)
	repository.SetPreloadKeys("User")
	return &authRepository{repository}
}

func (r *authRepository) FindByUserId(ctx context.Context, userId string) (*authmodel.Auth, error) {
	query := common.NewGormQuery()
	query.SetQuery("user_id = ?", userId)
	return r.FindOneByCond(ctx, query)
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*authmodel.Auth, error) {
	query := common.NewGormQuery()
	query.SetQuery("email = ?", email)
	return r.FindOneByCond(ctx, query)
}

func (r *authRepository) FindByPhone(ctx context.Context, phone string) (*authmodel.Auth, error) {
	query := common.NewGormQuery()
	query.SetQuery("phone = ?", phone)
	return r.FindOneByCond(ctx, query)
}

func (r *authRepository) FindByUserName(ctx context.Context, username string) (*authmodel.Auth, error) {
	query := common.NewGormQuery()
	query.SetQuery("email = ? OR phone = ?", username, username)
	return r.FindOneByCond(ctx, query)
}
