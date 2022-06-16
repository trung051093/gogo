package usermodel

import (
	"fmt"
	"log"
	"time"
	"user_management/common"
	"user_management/components/appctx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCreate struct {
	common.SQLModel
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
	if appCtx, ok := appctx.FromContext(ctx); ok {
		log.Println("🚀 ~ file: user.go ~ line 53 ~ ifappCtx,ok:=appctx.FromContext ~ ok", ok)

		rabbitmqService := appCtx.GetRabbitMQService()
		socketService := appCtx.GetSocketService()

		go func(user *UserCreate) {
			defer common.Recovery()
			dataIndex := &common.DataIndex{
				Id:       fmt.Sprintf("%d", user.Id),
				Index:    u.TableIndex(),
				Action:   common.Create,
				Data:     common.CompactJson(user),
				SendTime: time.Now(),
			}
			if publishErr := rabbitmqService.Publish(ctx, common.JsonToByte(dataIndex), []string{common.IndexingQueue}); publishErr != nil {
				log.Println("AfterCreate publish error:", publishErr)
			}
		}(u)

		go func(user *UserCreate) {
			defer common.Recovery()
			data := &common.Notification{
				Id: uuid.New().String(),
				Data: common.CompactJson(map[string]interface{}{
					"id": user.Id,
				}),
				Event:       fmt.Sprintf("%s-%s", user.TableName(), common.Create),
				Message:     "New user has been created",
				CreatedTime: time.Now(),
			}
			socketService.BroadcastToRoom(common.NotificationRoom, common.NotificationEvent, data)
		}(u)
	}
	return
}
