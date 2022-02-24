package common

import "time"

type SQLModel struct {
	Id        uint      `json:"id" gorm:"column:id;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoUpdateTime;"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime;"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deletedAt;"`
}
