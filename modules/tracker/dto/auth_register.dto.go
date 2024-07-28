package dto

import "gogo/modules/tracker/entity"

type AuthRegisterDto struct {
	Username string `json:"username" validate:"required" `
	Phone    string `json:"phone" validate:"required"`
	Password string `validate:"required" json:"password"`
}

func (dt *AuthRegisterDto) Validate() error {
	return nil
}

func (dt *AuthRegisterDto) ToEntity() *entity.Auth {
	return &entity.Auth{
		Username: dt.Username,
	}
}
