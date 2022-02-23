package user

import (
	usermodel "user_management/modules/user/model"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(newUser *usermodel.UserCreate) error {
	if err := s.repo.Create(newUser); err != nil {
		return err
	}
	return nil
}
