package usermodel

import (
	"log"
	"time"
	"user_management/common"
	"user_management/components/rabbitmq"

	"gorm.io/gorm"
)

type UserCreate struct {
	Id           int       `json:"id" gorm:"column:id;"`
	FirstName    string    `validate:"required" json:"firstName" gorm:"column:first_name;"`
	LastName     string    `validate:"required" json:"lastName" gorm:"column:last_name;"`
	Email        string    `validate:"required,email" json:"email" gorm:"column:email;"`
	Address      string    `json:"address" gorm:"column:address;"`
	Company      string    `json:"company" gorm:"column:company;"`
	BirthDate    time.Time `json:"birthDate" gorm:"column:birth_date;"`
	PhoneNumber  string    `json:"phoneNumber" gorm:"column:phone_number;"`
	Gender       string    `json:"gender" gorm:"column:gender;"`
	Role         string    `json:"role" gorm:"column:role;"`
	Password     string    `json:"-" gorm:"column:password;"`
	PasswordSalt string    `json:"-" gorm:"column:password_salt;"`
}

func (UserCreate) TableName() string { return User{}.TableName() }

func (u *UserCreate) AfterCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if rabbitmqService, ok := rabbitmq.FromContext(ctx); ok {
		go func() {
			defer common.Recovery()
			if publishErr := rabbitmqService.PublishWithTopic(ctx, common.IndexingQueue, u); publishErr != nil {
				log.Println("AfterCreate publish error:", publishErr)
			}
		}()
	}
	return
}
