package entity

import (
	"gogo/common"

	"github.com/google/uuid"
)

type GameHistory struct {
	common.SQLModel
	GameId  uuid.UUID `gorm:"column:game_id;" json:"game_id"`
	UserId  uuid.UUID `gorm:"column:user_id;" json:"user_id"`
	Game    *Game     `gorm:"foreignKey:GameId;references:Id" json:"game"`
	User    *User     `gorm:"foreignKey:UserId;references:Id" json:"user"`
	BuyIn   int       `gorm:"column:buy_in;" json:"buy_in"`
	CashOut int       `gorm:"column:cash_out;" json:"cash_out"`
	Total   int       `gorm:"column:total;" json:"total"`
}

func (GameHistory) EntityName() string { return "game_history" }

func (GameHistory) TableName() string { return "game_history" }
