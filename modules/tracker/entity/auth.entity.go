package entity

import (
	"gogo/common"

	"github.com/google/uuid"
)

type Auth struct {
	common.SQLModel
	Username     string    `gorm:"column:username;" json:"username"`
	UserId       uuid.UUID `gorm:"column:user_id;" json:"user_id" `
	Password     string    `gorm:"column:password;" json:"-" `
	PasswordSalt string    `gorm:"column:password_salt;" json:"-" `
	User         *User     `gorm:"foreignKey:UserId;references:Id" json:"user"`
}

func (Auth) EntityName() string { return "auth" }

func (Auth) TableName() string { return "auth" }
