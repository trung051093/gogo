package authmodel

import (
	"github.com/mitchellh/mapstructure"
)

type AuthCreateDto struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

func (dt *AuthCreateDto) Validate() error {
	return nil
}

func (dt *AuthCreateDto) ToEntity() *Auth {
	var auth Auth
	err := mapstructure.Decode(dt, &auth)
	if err != nil {
		return nil
	}

	return &auth
}
