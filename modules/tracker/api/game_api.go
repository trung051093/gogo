package api

import (
	"gogo/common"
	"gogo/components/appctx"
	"gogo/modules/tracker/dto"
	"gogo/modules/tracker/entity"
	"gogo/modules/tracker/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListGame godoc
//
//	@Summary		Get list of game
//	@Description	Get list of game
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			before	query		string																				false	"before cursor"
//	@Param			after	query		string																				false	"after cursor"
//	@Param			limit	query		int																					true	"limit"
//	@Success		200		{object}	common.ResponsePagination{data=[]entity.Game,pagination=common.CursorPagination}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games [get]
func ListGameHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		searchDto := &dto.GameSearchDto{
			CursorPagination: common.NewCursorPagination(nil, nil),
		}
		err := common.ParseRequest[dto.GameSearchDto](appCtx, ginCtx)(searchDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Game{}.EntityName(), err))

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		data, filter, paging, err := gameService.SearchCursorPaging(ginCtx.Request.Context(), searchDto)
		common.PanicIf(err != nil, common.ErrorCannotListEntity(entity.Game{}.EntityName(), err))

		ginCtx.JSON(http.StatusOK, common.SuccessResponsePagination(data, paging, filter))
	}
}

// GetGameByID godoc
//
//	@Summary		GetGameByID
//	@Description	GetGameByID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			before	query		string										false	"before cursor"
//	@Param			after	query		string										false	"after cursor"
//	@Param			limit	query		int											true	"limit"
//	@Success		200		{object}	common.ResponsePagination{data=entity.Game}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games/{game_id} [get]
func GetGameByIDHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id := ginCtx.Param("game_id")

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		game, err := gameService.FindById(ginCtx.Request.Context(), id)
		common.PanicIf(err != nil, common.ErrorCannotListEntity(entity.Game{}.EntityName(), err))

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(game))
	}
}

// NewGameHandler godoc
//
//	@Summary		NewGameHandler
//	@Description	NewGameHandler
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.GameCreateDto					true	"new"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games [post]
func NewGameHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		gameCreateDto := &dto.GameCreateDto{}

		err := common.ParseRequest[dto.GameCreateDto](appCtx, ginCtx)(gameCreateDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		game, err := gameService.New(ginCtx.Request.Context(), gameCreateDto)
		common.PanicIf(err != nil, err)

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse(game.Id.String()))
	}
}

// AddUserGame godoc
//
//	@Summary		AddUserGame
//	@Description	AddUserGame
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.GameAddUserDto					true	"buy in"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games/{game_id}/users [post]
func AddUserGameHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		addUserDto := &dto.GameAddUserDto{}

		err := common.ParseRequest[dto.GameAddUserDto](appCtx, ginCtx)(addUserDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		_, err = gameService.AddUser(ginCtx.Request.Context(), addUserDto)
		common.PanicIf(err != nil, err)

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse("OK"))
	}
}

// RemoveUser godoc
//
//	@Summary		RemoveUser
//	@Description	RemoveUser
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.GameRemoveUserDto					true	"buy in"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games/{game_id}/users/{user_id} [delete]
func RemoveUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		gameRemoveUserDto := &dto.GameRemoveUserDto{}

		err := common.ParseRequest[dto.GameRemoveUserDto](appCtx, ginCtx)(gameRemoveUserDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		_, err = gameService.RemoveUser(ginCtx.Request.Context(), gameRemoveUserDto)
		common.PanicIf(err != nil, err)

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse("OK"))
	}
}

// BuyinGameHandler godoc
//
//	@Summary		BuyinGameHandler
//	@Description	BuyinGameHandler
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.GameBuyInDto					true	"buy in"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games/{game_id}/users/{user_id}/buyin [post]
func BuyinGameHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		buyinGameDto := &dto.GameBuyInDto{}

		err := common.ParseRequest[dto.GameBuyInDto](appCtx, ginCtx)(buyinGameDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		_, err = gameService.Buyin(ginCtx.Request.Context(), buyinGameDto)
		common.PanicIf(err != nil, err)

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse("OK"))
	}
}

// CashoutGame godoc
//
//	@Summary		CashoutGame
//	@Description	CashoutGame
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.GameUserCashoutDto					true	"buy in"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games/{game_id}/users/{user_id}/cashout [post]
func CashoutGameHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		cashoutGameDto := &dto.GameUserCashoutDto{}

		err := common.ParseRequest[dto.GameUserCashoutDto](appCtx, ginCtx)(cashoutGameDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		_, err = gameService.Cashout(ginCtx.Request.Context(), cashoutGameDto)
		common.PanicIf(err != nil, err)

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse("OK"))
	}
}

// Summary godoc
//
//	@Summary		Summary
//	@Description	Summary
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.GameCashoutDto					true	"buy in"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/games/{game_id}/summary [post]
func SummaryHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		summaryGameDto := &dto.GameCashoutDto{}

		err := common.ParseRequest[dto.GameCashoutDto](appCtx, ginCtx)(summaryGameDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		gameService := service.NewGameServiceWithAppCtx(appCtx)
		_, err = gameService.Summary(ginCtx.Request.Context(), summaryGameDto)
		common.PanicIf(err != nil, err)

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse("OK"))
	}
}
