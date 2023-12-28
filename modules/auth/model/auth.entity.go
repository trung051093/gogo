package authmodel

import (
	"gogo/common"
	usermodel "gogo/modules/user/model"
)

const EntityName = "auth"

type Auth struct {
	common.SQLModel
	Phone        string          `json:"phone,omitempty" gorm:"column:phone;"`
	Email        string          `json:"email,omitempty" gorm:"column:email;"`
	UserId       string          `json:"userId" gorm:"column:user_id;"`
	Password     string          `json:"-" gorm:"column:password;"`
	PasswordSalt string          `json:"-" gorm:"column:password_salt;"`
	User         *usermodel.User `json:"user,omitempty" gorm:"foreignKey:UserId;references:Id"`
}

func (Auth) EntityName() string { return "auth" }

func (Auth) TableName() string { return "auth" }
