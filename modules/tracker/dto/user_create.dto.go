package dto

import (
	"gogo/modules/tracker/entity"

	"github.com/google/uuid"
)

type UserCreateDto struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (dt *UserCreateDto) Validate() error {
	return nil
}

func (dt *UserCreateDto) ToEntity() *entity.User {
	user := &entity.User{
		Name:  dt.Name,
		Phone: dt.Phone,
		Role:  entity.UserRole,
	}

	if dt.Id != "" {
		user.Id = uuid.MustParse(dt.Id)
	}

	return user
}
