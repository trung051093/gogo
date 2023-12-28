package usermodel

import (
	"fmt"
	"gogo/common"
	"gogo/components/appctx"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	common.SQLModel
	FirstName   string     `json:"firstName,omitempty" gorm:"column:first_name;"`
	LastName    string     `json:"lastName,omitempty" gorm:"column:last_name;"`
	Email       string     `json:"email,omitempty" gorm:"column:email;"`
	Address     string     `json:"address,omitempty" gorm:"column:address;"`
	Company     string     `json:"company,omitempty" gorm:"column:company;"`
	BirthDate   *time.Time `json:"birthDate,omitempty" gorm:"column:birth_date;"`
	Gender      string     `json:"gender,omitempty" gorm:"column:gender;"`
	PhoneNumber string     `json:"phoneNumber,omitempty" gorm:"column:phone_number;"`
	Role        string     `json:"role,omitempty" gorm:"column:role;"`
}

func (User) EntityName() string { return "user" }

func (User) TableName() string { return "user" }

func (User) TableIndex() string {
	return strings.ToLower(fmt.Sprintf("%s.%s", common.Project, User{}.TableName()))
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if appCtx, ok := appctx.FromContext(ctx); ok {
		rabbitmqService := appCtx.GetRabbitMQService()
		socketService := appCtx.GetSocketService()

		go func(user *User) {
			defer common.Recovery()
			dataIndex := &common.DataIndex{
				Id:       user.Id.String(),
				Index:    u.TableIndex(),
				Action:   common.Create,
				Data:     common.CompactJson(user),
				SendTime: time.Now(),
			}
			if publishErr := rabbitmqService.Publish(ctx, common.JsonToByte(dataIndex), []string{common.IndexingQueue}); publishErr != nil {
				log.Println("AfterCreate publish error:", publishErr)
			}
		}(u)

		go func(user *User) {
			defer common.Recovery()
			data := &common.Notification{
				Id: user.Id.String(),
				Data: map[string]interface{}{
					"id": user.Id,
				},
				Event:       fmt.Sprintf("%s-%s", user.TableName(), common.Create),
				Message:     "New user has been created",
				CreatedTime: time.Now(),
			}
			socketService.BroadcastToRoom(common.NotificationRoom, common.NotificationEvent, data)
		}(u)
	}
	return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context

	if appCtx, ok := appctx.FromContext(ctx); ok {
		rabbitmqService := appCtx.GetRabbitMQService()
		go func(user *User) {
			defer common.Recovery()
			dataIndex := &common.DataIndex{
				Id:       user.Id.String(),
				Index:    u.TableIndex(),
				Action:   common.Update,
				Data:     user,
				SendTime: time.Now(),
			}
			if publishErr := rabbitmqService.Publish(ctx, common.JsonToByte(dataIndex), []string{common.IndexingQueue}); publishErr != nil {
				log.Println("AfterUpdate publish error:", publishErr)
			}
		}(u)
	}
	return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context

	if appCtx, ok := appctx.FromContext(ctx); ok {
		rabbitmqService := appCtx.GetRabbitMQService()
		go func(user *User) {
			defer common.Recovery()
			dataIndex := &common.DataIndex{
				Id:       user.Id.String(),
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
