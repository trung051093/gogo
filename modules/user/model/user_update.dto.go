package usermodel

import (
	"time"

	"github.com/mitchellh/mapstructure"
)

type UserUpdateDto struct {
	FirstName   string     `json:"firstName"  validate:"omitempty" `
	LastName    string     `json:"lastName"  validate:"omitempty" `
	Email       string     `json:"email" validate:"omitempty,email"`
	Address     string     `json:"address"`
	Company     string     `json:"company"`
	BirthDate   *time.Time `json:"birthDate"`
	PhoneNumber string     `json:"phoneNumber"`
	Gender      string     `json:"gender"`
	Role        string     `json:"role"`
}

func (dt *UserUpdateDto) Validate() error {
	return nil
}

func (dt *UserUpdateDto) ToEntity(currentUser *User) *User {
	user := &User{}
	user.Id = currentUser.Id
	err := mapstructure.Decode(currentUser, user)
	if err != nil {
		return nil
	}

	return user
}

func (dt *UserUpdateDto) ToMapInterface() map[string]interface{} {
	updateDto := make(map[string]interface{})
	err := mapstructure.Decode(dt, updateDto)
	if err != nil {
		return nil
	}
	return updateDto
}
