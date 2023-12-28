package authmodel

import (
	usermodel "gogo/modules/user/model"
)

type AuthRegisterDto struct {
	Email     string `validate:"required,email" json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `validate:"required" json:"password"`
}

func (dt *AuthRegisterDto) Validate() error {
	return nil
}

func (dt *AuthRegisterDto) ToCreateUserDto() *usermodel.UserCreateDto {
	return &usermodel.UserCreateDto{
		Email:     dt.Email,
		FirstName: dt.FirstName,
		LastName:  dt.LastName,
		Role:      usermodel.UserRole,
	}
}

func (dt *AuthRegisterDto) ToCreateAuthDto() *AuthCreateDto {
	return &AuthCreateDto{
		Email:     dt.Email,
		FirstName: dt.FirstName,
		LastName:  dt.LastName,
		Password:  dt.Password,
	}
}
