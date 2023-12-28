package usermodel

import (
	"time"

	"github.com/mitchellh/mapstructure"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

func (s Role) String() string {
	switch s {
	case AdminRole:
		return "admin"
	case UserRole:
		return "user"
	}
	return "user"
}

type UserCreateDto struct {
	FirstName   string     `json:"firstName" validate:"required"`
	LastName    string     `json:"lastName" validate:"required" `
	Email       string     `json:"email" validate:"required,email"`
	Address     string     `json:"address"`
	Company     string     `json:"company"`
	BirthDate   *time.Time `json:"birthDate"`
	PhoneNumber string     `json:"phoneNumber"`
	Gender      string     `json:"gender"`
	Role        Role       `json:"role"`
	Password    string     `json:"password"`
}

func (dt *UserCreateDto) Validate() error {
	return nil
}

func (dt *UserCreateDto) ToEntity() *User {
	var user User
	err := mapstructure.Decode(dt, &user)
	if err != nil {
		return nil
	}

	return &user
}
