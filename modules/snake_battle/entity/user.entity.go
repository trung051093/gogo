package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a player in the game.
// Uses UUID for passwordless identification.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`
	Elo       int       `gorm:"default:1000;not null" json:"elo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relationships
	RoomMembers []RoomMember `gorm:"foreignKey:UserID" json:"roomMembers"` // A user can be a member of multiple rooms over time (though typically one active)
	GamesWon    []Game       `gorm:"foreignKey:WinnerID" json:"gamesWon"`  // Games where this user was the winner
	GamesLost   []Game       `gorm:"foreignKey:LoserID" json:"gamesLost"`  // Games where this user was the loser
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
