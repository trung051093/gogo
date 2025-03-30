package repository

import (
	"gogo/common"
	"gogo/modules/snake_battle/entity"

	"gorm.io/gorm"
)

type RoomMemberRepository interface {
	common.Repository[entity.RoomMember]
	CountMembersInRoom(roomID int) (int64, error)
	FindByUserAndRoom(userID string, roomID int) (*entity.RoomMember, error)
	// Add other specific methods, e.g., finding members by room, removing members
}

type roomMemberRepository struct {
	common.Repository[entity.RoomMember]
}

func NewRoomMemberRepository(db *gorm.DB) RoomMemberRepository {
	repo := common.NewRepository[entity.RoomMember](db)
	// Add preloads if needed, e.g., for User or Room details
	repo.SetPreloadKeys("User", "Room")
	return &roomMemberRepository{repo}
}

// Implement specific methods
func (r *roomMemberRepository) CountMembersInRoom(roomID int) (int64, error) {
	var count int64
	err := r.GetDB().Model(&entity.RoomMember{}).Where("room_id = ?", roomID).Count(&count).Error
	if err != nil {
		return 0, common.ErrorDB(err)
	}
	return count, nil
}

func (r *roomMemberRepository) FindByUserAndRoom(userID string, roomID int) (*entity.RoomMember, error) {
	var member entity.RoomMember
	err := r.GetDB().Where("user_id = ? AND room_id = ?", userID, roomID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrorCannotFoundEntity(entity.RoomMember{}.TableName(), err)
		}
		return nil, common.ErrorDB(err)
	}
	return &member, nil
}
