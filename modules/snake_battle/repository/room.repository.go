package repository

import (
	"context"
	"gogo/common"
	"gogo/modules/snake_battle/entity"

	"gorm.io/gorm"
)

type RoomRepository interface {
	common.Repository[entity.Room]
	// FindByCode finds a room by its unique code.
	FindByCode(ctx context.Context, roomCode string) (*entity.Room, error)
	// Add other specific room methods if needed
}

type roomRepository struct {
	common.Repository[entity.Room]
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	repo := common.NewRepository[entity.Room](db)
	// Set preloads if needed, e.g., for Members
	// repo.SetPreloadKeys("Members")
	return &roomRepository{repo}
}

// Implement specific methods
func (r *roomRepository) FindByCode(ctx context.Context, roomCode string) (*entity.Room, error) {
	var room entity.Room
	// Use GetDBWithContext if available in common.Repository, otherwise GetDB()
	db := r.GetDB().WithContext(ctx) // Assuming GetDB() exists and we add context
	err := db.Where("room_code = ?", roomCode).First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrorCannotFoundEntity(entity.Room{}.TableName(), err)
		}
		return nil, common.ErrorDB(err)
	}
	return &room, nil
}
