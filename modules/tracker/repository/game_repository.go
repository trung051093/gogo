package repository

import (
	"gogo/common"
	"gogo/modules/tracker/entity"

	"gorm.io/gorm"
)

type GameRepository interface {
	common.Repository[entity.Game]
}

type gameRepository struct {
	common.Repository[entity.Game]
}

func NewGameRepository(db *gorm.DB) GameRepository {
	repository := common.NewRepository[entity.Game](db)
	repository.SetPreloadKeys("GameHistory", "GameHistory.User")
	return &gameRepository{repository}
}
