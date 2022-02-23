package usermodel

import (
	"time"
	"user_management/common"
)

type User struct {
	common.SQLModel
	FirstName   string    `json:"firstName" gorm:"column:first_name;"`
	LastName    string    `json:"lastName" gorm:"column:last_name;"`
	Email       string    `json:"email" gorm:"column:email;"`
	Address     string    `json:"address" gorm:"column:address;"`
	Company     string    `json:"company" gorm:"column:company;"`
	BirthDate   time.Time `json:"birthDate" gorm:"column:birth_date;"`
	PhoneNumber string    `json:"phoneNumber" gorm:"column:phone_number;"`
}

func (User) TableName() string { return "users" }
