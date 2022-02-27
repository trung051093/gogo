package user

import (
	"user_management/common"
	usermodel "user_management/modules/user/model"
)

type userService struct {
	repo *userRepository
}

func NewUserService(repo *userRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) SearchUsers(cond map[string]interface{}, p *common.Pagination) ([]usermodel.User, error) {
	return s.repo.Search(cond, p)
}

func (s *userService) GetUser(id uint) (*usermodel.User, error) {
	return s.repo.Get(id)
}

func (s *userService) CreateUser(newUser *usermodel.UserCreate) error {
	return s.repo.Create(newUser)
}

func (s *userService) UpdateUser(id uint, userUpdate *usermodel.UserUpdate) error {
	user, err := s.GetUser(id)
	if user == nil {
		return common.ErrNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.Update(map[string]interface{}{"id": id}, userUpdate)
}

func (s *userService) DeleteUser(id uint) error {
	user, err := s.GetUser(id)
	if user == nil {
		return common.ErrNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.Delete(map[string]interface{}{"id": id})
}
