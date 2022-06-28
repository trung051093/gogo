package common

import (
	"time"
)

type SQLModel struct {
	Id        int        `json:"id" gorm:"column:id;"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime;"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime;"`
	IsActive  int        `json:"isActive" gorm:"column:is_active;default:1;"`
}
