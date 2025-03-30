package entity

import (
	"time"

	"github.com/google/uuid"
)

// RoomMember links a User to a Room.
type RoomMember struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"userId"` // Foreign key to User
	RoomID   uint      `gorm:"not null" json:"roomId"`           // Foreign key to Room
	JoinedAt time.Time `json:"joinedAt"`
	IsHost   bool      `gorm:"default:false" json:"isHost"`

	// Relationships (Belongs To)
	User User `gorm:"foreignKey:UserID" json:"user"`
	Room Room `gorm:"foreignKey:RoomID" json:"room"`
}

// TableName specifies the table name for GORM
func (RoomMember) TableName() string {
	return "room_members"
}
