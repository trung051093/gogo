package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLModel struct {
	Id        uuid.UUID  `json:"id" gorm:"column:id;primary_key;default:uuid_generate_v4()"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;"`
}

func (base *SQLModel) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("id", uuid)
	return nil
}
