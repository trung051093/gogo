package usermodel

import (
	"fmt"
	"log"
	"time"
	"user_management/common"
	rabbitmqprovider "user_management/components/rabbitmq"

	"gorm.io/gorm"
)

type UserUpdate struct {
	common.SQLModel
	FirstName    string     `validate:"omitempty" json:"firstName" gorm:"column:first_name;"`
	LastName     string     `validate:"omitempty" json:"lastName" gorm:"column:last_name;"`
	Email        *string    `validate:"omitempty,email" json:"email" gorm:"column:email;"`
	Address      *string    `json:"address" gorm:"column:address;"`
	Company      *string    `json:"company" gorm:"column:company;"`
	BirthDate    *time.Time `json:"birthDate" gorm:"column:birth_date;"`
	PhoneNumber  *string    `json:"phoneNumber" gorm:"column:phone_number;"`
	Gender       string     `json:"gender" gorm:"column:gender;"`
	Role         string     `json:"role" gorm:"column:role;"`
	Password     string     `json:"-" gorm:"column:password;"`
	PasswordSalt string     `json:"-" gorm:"column:password_salt;"`
}

func (UserUpdate) TableName() string  { return User{}.TableName() }
func (UserUpdate) TableIndex() string { return User{}.TableIndex() }

func (u *UserUpdate) AfterUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if rabbitmqService, ok := rabbitmqprovider.FromContext(ctx); ok {
		log.Println("AfterUpdate rabbitmqService:", rabbitmqService)
		go func(user *UserUpdate) {
			defer common.Recovery()
			dataIndex := &common.DataIndex{
				Id:       fmt.Sprintf("%d", user.Id),
				Index:    u.TableIndex(),
				Action:   common.Update,
				Data:     common.CompactJson(user),
				SendTime: time.Now(),
			}
			if publishErr := rabbitmqService.Publish(ctx, common.JsonToByte(dataIndex), []string{common.IndexingQueue}); publishErr != nil {
				log.Println("AfterUpdate publish error:", publishErr)
			}
		}(u)
	}
	return
}
