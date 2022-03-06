package usermodel

import (
	"time"
	"user_management/common"
)

const EntityName = "user"

type User struct {
	common.SQLModel
	FirstName   string    `json:"firstName,omitempty" gorm:"column:first_name;"`
	LastName    string    `json:"lastName,omitempty" gorm:"column:last_name;"`
	Email       string    `json:"email,omitempty" gorm:"column:email;"`
	Address     string    `json:"address,omitempty" gorm:"column:address;"`
	Company     string    `json:"company,omitempty" gorm:"column:company;"`
	BirthDate   time.Time `json:"birthDate,omitempty" gorm:"column:birth_date;"`
	Gender      string    `json:"gender,omitempty" gorm:"column:gender;"`
	PhoneNumber string    `json:"phoneNumber,omitempty" gorm:"column:phone_number;"`
}

var UserField = map[string]string{
	"id":          "id",
	"firstName":   "first_name",
	"lastName":    "last_name",
	"email":       "email",
	"address":     "address",
	"company":     "company",
	"birthDate":   "birth_date",
	"gender":      "gender",
	"phoneNumber": "phone_number",
}

func (User) TableName() string { return "users" }
