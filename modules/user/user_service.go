package user

import (
	"context"
	"user_management/common"
	esprovider "user_management/components/elasticsearch"
	elasticsearchmodel "user_management/components/elasticsearch/model"
	usermodel "user_management/modules/user/model"
)

type userService struct {
	repo      UserRepository
	esService esprovider.ElasticSearchSevice
}

func NewUserService(repo UserRepository, esService esprovider.ElasticSearchSevice) UserService {
	return &userService{
		repo:      repo,
		esService: esService,
	}
}

func (s *userService) SearchUsers(ctx context.Context, cond map[string]interface{}, f *usermodel.UserFilter, p *common.Pagination) ([]usermodel.User, error) {
	return s.repo.Search(ctx, cond, f, p)
}

func (s *userService) SearchUser(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error) {
	return s.repo.SearchOne(ctx, cond)
}

func (s *userService) GetUser(ctx context.Context, id int) (*usermodel.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, newUser *usermodel.UserCreate) (int, error) {
	return s.repo.Create(ctx, newUser)
}

func (s *userService) UpdateUser(ctx context.Context, id int, userUpdate *usermodel.UserUpdate) (int, error) {
	user, err := s.GetUser(ctx, id)
	if user == nil || err != nil {
		return -1, err
	}
	userUpdate.Id = user.Id
	return s.repo.Update(ctx, id, userUpdate)
}

func (s *userService) DeleteUser(ctx context.Context, id int) (int, error) {
	user, err := s.GetUser(ctx, id)
	if user == nil || err != nil {
		return -1, err
	}
	return s.repo.Delete(ctx, user)
}

func (s *userService) EsSearch(ctx context.Context, query string, lastIndex string, f *usermodel.UserFilter, p *common.Pagination) (*elasticsearchmodel.SearchResults, error) {
	userEsQuery := &usermodel.UserEsQuery{
		Query:     query,
		LastIndex: lastIndex,
		Paging:    p,
		Filter:    f,
	}
	esUserQuery := usermodel.GetUserESQuery(ctx, userEsQuery)
	return s.esService.Search(ctx, usermodel.User{}.TableIndex(), esUserQuery)
}
