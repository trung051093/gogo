package user

import (
	"context"
	"user_management/common"
	decorator "user_management/decorators"
	usermodel "user_management/modules/user/model"
)

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) SearchUsersTrace(ctx context.Context, cond map[string]interface{}, f *usermodel.UserFilter, p *common.Pagination) ([]usermodel.User, error) {
	data, err := decorator.TraceService(ctx, "userService.SearchUsers")(s, "SearchUsers")(ctx, cond, f, p)
	return data.([]usermodel.User), err
}

func (s *userService) SearchUsers(ctx context.Context, cond map[string]interface{}, f *usermodel.UserFilter, p *common.Pagination) ([]usermodel.User, error) {
	return s.repo.Search(ctx, cond, f, p)
}

func (s *userService) SearchUser(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error) {
	return s.repo.SearchOne(ctx, cond)
}

func (s *userService) GetUser(ctx context.Context, id uint) (*usermodel.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, newUser *usermodel.UserCreate) (int, error) {
	return s.repo.Create(ctx, newUser)
}

func (s *userService) UpdateUser(ctx context.Context, id uint, userUpdate *usermodel.UserUpdate) error {
	user, err := s.GetUser(ctx, id)
	if user == nil || err != nil {
		return err
	}
	userUpdate.Id = user.Id
	return s.repo.Update(ctx, id, userUpdate)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	user, err := s.GetUser(ctx, id)
	if user == nil || err != nil {
		return err
	}
	return s.repo.Delete(ctx, user)
}
