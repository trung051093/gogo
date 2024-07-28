package service

import (
	"context"
	"gogo/common"
	"gogo/components/appctx"
	"gogo/components/hasher"
	"gogo/modules/tracker/dto"
	"gogo/modules/tracker/entity"
	"gogo/modules/tracker/repository"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type GameService interface {
	common.Service[entity.Game, repository.GameRepository]
	New(ctx context.Context, payload *dto.GameCreateDto) (*entity.Game, error)
	Buyin(ctx context.Context, payload *dto.GameBuyInDto) (*entity.Game, error)
	Cashout(ctx context.Context, payload *dto.GameUserCashoutDto) (*entity.Game, error)
	Summary(ctx context.Context, payload *dto.GameCashoutDto) (*entity.Game, error)
	AddUser(ctx context.Context, payload *dto.GameAddUserDto) (*entity.Game, error)
	RemoveUser(ctx context.Context, payload *dto.GameRemoveUserDto) (*entity.Game, error)
}

type gameService struct {
	common.Service[entity.Game, repository.GameRepository]
	gameRepo        repository.GameRepository
	gameHistoryRepo repository.GameHistoryRepository
	authService     AuthService
	userService     UserService
}

func NewGameService(
	gameRepo repository.GameRepository,
	gameHistoryRepo repository.GameHistoryRepository,
	authService AuthService,
	userService UserService,
) GameService {
	return &gameService{
		Service:         common.NewService[entity.Game, repository.GameRepository](gameRepo),
		gameRepo:        gameRepo,
		gameHistoryRepo: gameHistoryRepo,
		authService:     authService,
		userService:     userService,
	}
}

func NewGameServiceWithAppCtx(appCtx appctx.AppContext) GameService {
	appConfig := appCtx.GetConfig()
	authRepo := repository.NewAuthRepository(appCtx.GetMainDBConnection())
	userRepo := repository.NewUserRepository(appCtx.GetMainDBConnection())
	gameRepo := repository.NewGameRepository(appCtx.GetMainDBConnection())
	gameHistoryRepo := repository.NewGameHistoryRepository(appCtx.GetMainDBConnection())

	cacheService := appCtx.GetCacheService()
	hashService := hasher.NewHashService()
	jwtProvider := NewJWTProvider(appConfig.JWT.Secret)

	userService := NewUserService(userRepo, appConfig)
	authService := NewAuthService(authRepo, userRepo, appConfig, jwtProvider, hashService, cacheService)
	gameService := NewGameService(gameRepo, gameHistoryRepo, authService, userService)

	return gameService
}

func (s *gameService) New(ctx context.Context, payload *dto.GameCreateDto) (*entity.Game, error) {
	game, err := s.gameRepo.Create(ctx, payload.ToEntity())
	common.PanicIf(err != nil || game == nil, common.ErrorCannotCreateEntity(entity.Game{}.EntityName(), err))

	userEntities := payload.ToUserEntities()
	newUsers := lo.Filter(userEntities, func(item *entity.User, _ int) bool {
		return item.Id == uuid.Nil
	})
	existUsers := lo.Filter(userEntities, func(item *entity.User, _ int) bool {
		return item.Id != uuid.Nil
	})
	users, err := s.userService.GetRepository().CreateList(ctx, newUsers)
	common.PanicIf(err != nil || users == nil, common.ErrorCannotCreateEntity(entity.User{}.EntityName(), err))

	gameHistories, err := s.gameHistoryRepo.AddUsers(ctx, game.Id, append(existUsers, users...), payload.DefaultBuyIn)
	common.PanicIf(err != nil || gameHistories == nil, common.ErrorCannotCreateEntity(entity.GameHistory{}.EntityName(), err))

	return game, nil
}

func (s *gameService) AddUser(ctx context.Context, payload *dto.GameAddUserDto) (*entity.Game, error) {
	game, err := s.gameRepo.FindById(ctx, payload.GameId)
	common.PanicIf(err != nil || game == nil, common.ErrorCannotFoundEntity(entity.Game{}.EntityName(), err))

	userEntities := payload.ToUserEntities()
	newUsers := lo.Filter(userEntities, func(item *entity.User, _ int) bool {
		return item.Id == uuid.Nil
	})
	existUsers := lo.Filter(userEntities, func(item *entity.User, _ int) bool {
		return item.Id != uuid.Nil
	})
	users, err := s.userService.GetRepository().CreateList(ctx, newUsers)
	common.PanicIf(err != nil || users == nil, common.ErrorCannotCreateEntity(entity.User{}.EntityName(), err))

	addUsers := append(existUsers, users...)
	gameUsers := lo.Map(game.GameHistory, func(item *entity.GameHistory, _ int) *entity.User {
		return item.User
	})
	filteredUsers := lo.Filter(addUsers, func(item *entity.User, _ int) bool {
		if _, isExistInGame := lo.Find(gameUsers, func(gameUser *entity.User) bool {
			return gameUser.Id == item.Id
		}); isExistInGame {
			return false
		}
		return true
	})

	gameHistories, err := s.gameHistoryRepo.AddUsers(ctx, game.Id, filteredUsers, payload.DefaultBuyIn)
	common.PanicIf(err != nil || gameHistories == nil, common.ErrorCannotCreateEntity(entity.GameHistory{}.EntityName(), err))

	err = s.RecalculateTotal(ctx, payload.GameId)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.Game{}.EntityName(), err))

	return game, nil
}

func (s *gameService) RemoveUser(ctx context.Context, payload *dto.GameRemoveUserDto) (*entity.Game, error) {
	game, err := s.gameRepo.FindById(ctx, payload.GameId)
	common.PanicIf(err != nil || game == nil, common.ErrorCannotFoundEntity(entity.Game{}.EntityName(), err))

	err = s.gameHistoryRepo.RemoveUser(ctx, game.Id, uuid.MustParse(payload.UserId))
	common.PanicIf(err != nil, common.ErrorCannotCreateEntity(entity.GameHistory{}.EntityName(), err))

	err = s.RecalculateTotal(ctx, payload.GameId)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.Game{}.EntityName(), err))

	return game, nil
}

func (s *gameService) Buyin(ctx context.Context, payload *dto.GameBuyInDto) (*entity.Game, error) {
	game, err := s.gameRepo.FindById(ctx, payload.GameId)
	common.PanicIf(err != nil || game == nil, common.ErrorCannotFoundEntity(entity.Game{}.EntityName(), err))

	gameHistory, err := s.gameHistoryRepo.FindByGameIdAndUserId(ctx, payload.GameId, payload.UserId)
	common.PanicIf(err != nil || gameHistory == nil, common.ErrorCannotFoundEntity(entity.GameHistory{}.EntityName(), err))

	gameHistory.BuyIn = gameHistory.BuyIn + payload.BuyIn
	_, err = s.gameHistoryRepo.Save(ctx, gameHistory)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.GameHistory{}.EntityName(), err))

	err = s.RecalculateTotal(ctx, payload.GameId)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.Game{}.EntityName(), err))

	return game, nil
}

func (s *gameService) Cashout(ctx context.Context, payload *dto.GameUserCashoutDto) (*entity.Game, error) {
	game, err := s.gameRepo.FindById(ctx, payload.GameId)
	common.PanicIf(err != nil || game == nil, common.ErrorCannotFoundEntity(entity.Game{}.EntityName(), err))

	gameHistory, err := s.gameHistoryRepo.FindByGameIdAndUserId(ctx, payload.GameId, payload.UserId)
	common.PanicIf(err != nil || gameHistory == nil, common.ErrorCannotFoundEntity(entity.GameHistory{}.EntityName(), err))

	gameHistory.CashOut = payload.Cashout
	gameHistory.Total = gameHistory.BuyIn + gameHistory.CashOut
	_, err = s.gameHistoryRepo.Save(ctx, gameHistory)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.GameHistory{}.EntityName(), err))

	err = s.RecalculateTotal(ctx, payload.GameId)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.Game{}.EntityName(), err))

	return game, nil
}

func (s *gameService) Summary(ctx context.Context, payload *dto.GameCashoutDto) (*entity.Game, error) {
	game, err := s.gameRepo.FindById(ctx, payload.GameId)
	common.PanicIf(err != nil || game == nil, common.ErrorCannotFoundEntity(entity.Game{}.EntityName(), err))

	gameHistories, err := s.gameHistoryRepo.FindByGameId(ctx, payload.GameId)
	common.PanicIf(err != nil, common.ErrorCannotFoundEntity(entity.GameHistory{}.EntityName(), err))

	userCashoutById := lo.KeyBy(payload.CashoutSummary, func(item *dto.UserCashoutDto) string {
		return item.UserId
	})

	for _, gameHistory := range gameHistories {
		userCashout := userCashoutById[gameHistory.UserId.String()]
		gameHistory.CashOut = userCashout.Cashout
		gameHistory.Total = gameHistory.BuyIn + gameHistory.CashOut
		_, err = s.gameHistoryRepo.Save(ctx, gameHistory)
		common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.GameHistory{}.EntityName(), err))
	}

	err = s.RecalculateTotal(ctx, payload.GameId)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.Game{}.EntityName(), err))

	return game, nil
}

func (s *gameService) RecalculateTotal(ctx context.Context, gameId string) error {
	game, err := s.gameRepo.FindById(ctx, gameId)
	common.PanicIf(err != nil || game == nil, common.ErrorCannotFoundEntity(entity.Game{}.EntityName(), err))

	totalBuyIn := lo.Reduce(game.GameHistory, func(total int, item *entity.GameHistory, _ int) int {
		return total + item.BuyIn
	}, 0)

	totalCashout := lo.Reduce(game.GameHistory, func(total int, item *entity.GameHistory, _ int) int {
		return total + item.CashOut
	}, 0)

	game.TotalCashOut = totalCashout
	game.TotalBuyIn = totalBuyIn
	_, err = s.gameRepo.Save(ctx, game)
	common.PanicIf(err != nil, common.ErrorCannotUpdateEntity(entity.Game{}.EntityName(), err))

	return nil
}
