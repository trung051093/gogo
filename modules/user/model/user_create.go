package usermodel

import (
	"fmt"
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

func (UserCreate) TableIndex() string { return User{}.TableIndex() }

func (u *UserCreate) AfterCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if rabbitmqService, ok := rabbitmq.FromContext(ctx); ok {
		log.Println("AfterCreate rabbitmqService:", rabbitmqService)
		go func() {
			defer common.Recovery()
			dataIndex := &common.DataIndex{
				Id:       fmt.Sprintf("%d", u.Id),
				Index:    u.TableIndex(),
				Action:   common.Create,
				Data:     common.CompactJson(u),
				SendTime: time.Now(),
			}
			if publishErr := rabbitmqService.PublishWithTopic(ctx, common.IndexingQueue, dataIndex); publishErr != nil {
				log.Println("AfterCreate publish error:", publishErr)
			}
		}()
	}
	return
}
