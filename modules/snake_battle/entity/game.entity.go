package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Game records the result of a match.
type Game struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoomID    uint           `gorm:"not null" json:"roomId"`    // Foreign key to Room
	WinnerID  *uuid.UUID     `gorm:"type:uuid" json:"winnerId"` // Foreign key to User (nullable)
	LoserID   *uuid.UUID     `gorm:"type:uuid" json:"loserId"`  // Foreign key to User (nullable)
	EloChange *int           `json:"eloChange"`                 // ELO change (nullable, e.g., for draws or if not tracked)
	StartedAt time.Time      `json:"startedAt"`
	EndedAt   *time.Time     `json:"endedAt"`                    // Nullable until the game finishes
	GameData  datatypes.JSON `gorm:"type:jsonb" json:"gameData"` // Optional game replay/state data

	// Relationships (Belongs To)
	Room   Room  `gorm:"foreignKey:RoomID" json:"room"`
	Winner *User `gorm:"foreignKey:WinnerID" json:"winner"` // Pointer because WinnerID is nullable
	Loser  *User `gorm:"foreignKey:LoserID" json:"loser"`   // Pointer because LoserID is nullable
}

// TableName specifies the table name for GORM
func (Game) TableName() string {
	return "games"
}
