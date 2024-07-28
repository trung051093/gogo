package entity

import (
	"gogo/common"
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
	return "unknown"
}

type User struct {
	common.SQLModel
	Name  string `gorm:"column:name;" json:"name"`
	Phone string `gorm:"column:phone;" json:"phone"`
	Role  Role   `gorm:"column:role;" json:"-"`
}

func (User) EntityName() string { return "users" }

func (User) TableName() string { return "users" }
