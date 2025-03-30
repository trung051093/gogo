package repository

import (
	"gogo/common"
	"gogo/modules/snake_battle/entity"

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
	repository.SetPreloadKeys("Room", "Winner", "Loser")
	return &gameRepository{repository}
}
