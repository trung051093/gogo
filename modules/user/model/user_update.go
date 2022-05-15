package usermodel

import (
	"time"
)

type UserUpdate struct {
	FirstName    string     `validate:"omitempty" json:"firstName" gorm:"column:first_name;"`
	LastName     string     `validate:"omitempty" json:"lastName" gorm:"column:last_name;"`
	Email        *string    `validate:"omitempty,email" json:"email" gorm:"column:email;"`
	Address      *string    `json:"address" gorm:"column:address;"`
	Company      *string    `json:"company" gorm:"column:company;"`
	BirthDate    *time.Time `json:"birthDate" gorm:"column:birth_date;"`
	PhoneNumber  *string    `json:"phoneNumber" gorm:"column:phone_number;"`
	Gender       string     `json:"gender" gorm:"column:gender;"`
	Role         string     `json:"role" gorm:"column:role;"`
	Password     string     `json:"password" gorm:"column:password;"`
	PasswordSalt string     `json:"-" gorm:"column:password_salt;"`
}

func (UserUpdate) TableName() string { return User{}.TableName() }
