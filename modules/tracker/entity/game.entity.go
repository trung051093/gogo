package entity

import (
	"gogo/common"
)

type Game struct {
	common.SQLModel
	Name         string         `gorm:"column:name;" json:"name"`
	TotalBuyIn   int            `gorm:"column:total_buy_in;" json:"total_buy_in"`
	TotalCashOut int            `gorm:"column:total_cash_out;" json:"total_cash_out"`
	GameHistory  []*GameHistory `gorm:"foreignKey:GameId" json:"history"`
}

func (Game) EntityName() string { return "game" }

func (Game) TableName() string { return "game" }
