package common

import "time"

type SQLModel struct {
	Id        *int       `json:"id" gorm:"column:id;"`
	CreatedAt *time.Time `json:"createdAt,omitempty" gorm:"column:created_at;autoUpdateTime;"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at;autoUpdateTime;"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"column:deletedAt;"`
}
