package repository

import (
	"context"
	"gogo/common"
	"gogo/modules/tracker/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameHistoryRepository interface {
	common.Repository[entity.GameHistory]
	FindByGameId(ctx context.Context, gameId string) ([]*entity.GameHistory, error)
	FindByGameIdAndUserId(ctx context.Context, gameId string, userId string) (*entity.GameHistory, error)
	AddUsers(ctx context.Context, gameId uuid.UUID, users []*entity.User, defaultBuyIn int) ([]*entity.GameHistory, error)
	RemoveUser(ctx context.Context, gameId uuid.UUID, userId uuid.UUID) error
}

type gameHistoryRepository struct {
	common.Repository[entity.GameHistory]
}

func NewGameHistoryRepository(db *gorm.DB) GameHistoryRepository {
	repository := common.NewRepository[entity.GameHistory](db)
	return &gameHistoryRepository{repository}
}

func (r *gameHistoryRepository) AddUsers(ctx context.Context, gameId uuid.UUID, users []*entity.User, defaultBuyIn int) ([]*entity.GameHistory, error) {
	var gameHistories []*entity.GameHistory
	for _, user := range users {
		gameHistory := &entity.GameHistory{
			GameId:  gameId,
			UserId:  user.Id,
			BuyIn:   defaultBuyIn,
			CashOut: 0,
		}
		gameHistories = append(gameHistories, gameHistory)
	}

	return r.CreateList(ctx, gameHistories)
}

func (r *gameHistoryRepository) RemoveUser(ctx context.Context, gameId uuid.UUID, userId uuid.UUID) error {
	return r.GetDB().Where("game_id = ? AND user_id = ?", gameId, userId).Delete(&entity.GameHistory{}).Error
}

func (r *gameHistoryRepository) FindByGameId(ctx context.Context, gameId string) ([]*entity.GameHistory, error) {
	query := common.NewGormQuery()
	query.SetQuery("game_id = ?", gameId)
	return r.FindByCond(ctx, query)
}

func (r *gameHistoryRepository) FindByGameIdAndUserId(ctx context.Context, gameId string, userId string) (*entity.GameHistory, error) {
	query := common.NewGormQuery()
	query.SetQuery("game_id = ?", gameId)
	query.SetQuery("user_id = ?", userId)
	return r.FindOneByCond(ctx, query)
}
