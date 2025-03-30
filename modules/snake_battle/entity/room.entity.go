package entity

import (
	"time"
)

// Room represents a game room.
// Status column has been removed.
type Room struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoomCode   string    `gorm:"type:varchar(10);unique;not null" json:"roomCode"`
	MaxPlayers int       `gorm:"default:2;not null" json:"maxPlayers"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	// Relationships
	Members []RoomMember `gorm:"foreignKey:RoomID" json:"members"` // Members currently in the room
	Games   []Game       `gorm:"foreignKey:RoomID" json:"games"`   // Games played in this room
}

// TableName specifies the table name for GORM
func (Room) TableName() string {
	return "rooms"
}
