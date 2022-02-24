package user

import (
	"user_management/common"
	usermodel "user_management/modules/user/model"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SearchUsers(cond map[string]interface{}, p *common.Pagination) ([]usermodel.User, error) {
	return s.repo.Search(cond, p)
}

func (s *UserService) GetUser(id uint) (*usermodel.User, error) {
	return s.repo.Get(id)
}

func (s *UserService) CreateUser(newUser *usermodel.UserCreate) error {
	return s.repo.Create(newUser)
}

func (s *UserService) UpdateUser(id uint, userUpdate *usermodel.UserUpdate) error {
	user, err := s.GetUser(id)
	if user == nil {
		return common.ErrNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.Update(map[string]interface{}{"id": id}, userUpdate)
}

func (s *UserService) DeleteUser(id uint) error {
	user, err := s.GetUser(id)
	if user == nil {
		return common.ErrNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.Delete(map[string]interface{}{"id": id})
}
