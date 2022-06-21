package usermodel

import (
	"fmt"
	"log"
	"strings"
	"time"
	"user_management/common"
	"user_management/components/appctx"

	"gorm.io/gorm"
)

const EntityName = "user"

type User struct {
	common.SQLModel
	FirstName    string     `json:"firstName" gorm:"column:first_name;"`
	LastName     string     `json:"lastName" gorm:"column:last_name;"`
	Email        string     `json:"email" gorm:"column:email;"`
	Address      string     `json:"address" gorm:"column:address;"`
	Company      string     `json:"company" gorm:"column:company;"`
	BirthDate    *time.Time `json:"birthDate" gorm:"column:birth_date;"`
	Gender       string     `json:"gender" gorm:"column:gender;"`
	PhoneNumber  string     `json:"phoneNumber" gorm:"column:phone_number;"`
	Role         string     `json:"role" gorm:"column:role;default:user"`
	Password     string     `json:"-" gorm:"column:password;"`
	PasswordSalt string     `json:"-" gorm:"column:password_salt;"`
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
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

func (User) TableName() string { return "users" }

func (User) TableIndex() string {
	return strings.ToLower(fmt.Sprintf("%s-%s", common.Project, User{}.TableName()))
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context

	if appCtx, ok := appctx.FromContext(ctx); ok {
		rabbitmqService := appCtx.GetRabbitMQService()
		go func(user *User) {
			defer common.Recovery()
			dataIndex := &common.DataIndex{
				Id:       fmt.Sprintf("%d", user.Id),
				Index:    u.TableIndex(),
				Action:   common.Delete,
				Data:     user,
				SendTime: time.Now(),
			}
			if publishErr := rabbitmqService.Publish(ctx, common.JsonToByte(dataIndex), []string{common.IndexingQueue}); publishErr != nil {
				log.Println("AfterDelete publish error:", publishErr)
			}
		}(u)
	}
	return
}
