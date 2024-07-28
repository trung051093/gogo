package service

import (
	"gogo/common"
	"gogo/components/appctx"
	"gogo/modules/tracker/entity"
	"gogo/modules/tracker/repository"
)

type UserService interface {
	common.Service[entity.User, repository.UserRepository]
}

type userService struct {
	common.Service[entity.User, repository.UserRepository]
	appConfig *appctx.Config
}

func NewUserService(
	repo repository.UserRepository,
	appConfig *appctx.Config,
) UserService {
	service := common.NewService(repo)
	return &userService{service, appConfig}
}

func NewUserServiceWithAppCtx(appCtx appctx.AppContext) UserService {
	userRepo := repository.NewUserRepository(appCtx.GetMainDBConnection())
	appConfig := appCtx.GetConfig()
	userService := NewUserService(userRepo, appConfig)
	return userService
}
