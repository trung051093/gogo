package user

import (
	"context"
	"gogo/common"
	"gogo/components/appctx"
	esprovider "gogo/components/elasticsearch"
	elasticsearchmodel "gogo/components/elasticsearch/model"
	usermodel "gogo/modules/user/model"
)

type UserService interface {
	common.Service[usermodel.User, UserRepository]
	EsSearch(ctx context.Context, userEsQuery *usermodel.UserEsSearchDto) (*elasticsearchmodel.SearchResults, error)
}

type userService struct {
	common.Service[usermodel.User, UserRepository]
	appConfig *appctx.Config
	esService esprovider.ElasticSearchSevice
}

func NewUserService(
	repo UserRepository,
	appConfig *appctx.Config,
	esService esprovider.ElasticSearchSevice,
) UserService {
	service := common.NewService[usermodel.User, UserRepository](repo)
	return &userService{service, appConfig, esService}
}

func NewUserServiceWithAppCtx(appCtx appctx.AppContext) UserService {
	userRepo := NewUserRepository(appCtx.GetMainDBConnection())
	appConfig := appCtx.GetConfig()
	esService := appCtx.GetESService()
	userService := NewUserService(userRepo, appConfig, esService)
	return userService
}

func (s *userService) EsSearch(ctx context.Context, userEsQuery *usermodel.UserEsSearchDto) (*elasticsearchmodel.SearchResults, error) {
	return s.esService.Search(ctx, usermodel.User{}.TableIndex(), userEsQuery.ToQuery())
}
