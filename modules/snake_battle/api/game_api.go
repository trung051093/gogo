package api

import (
	"gogo/common"
	"gogo/components/appctx"
	"gogo/modules/snake_battle/dto"
	"gogo/modules/snake_battle/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
//
//	@Summary		Register
//	@Description	Register
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.CreateRoomReq		true		"create room"
//	@Success		200		{object}	dto.Response{data=CreateRoomRes}	"response"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/room [post]
func CreateRoom(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		createRoomReq := &dto.CreateRoomReq{}
		if err := common.ParseRequest[dto.CreateRoomReq](appCtx, ginCtx)(createRoomReq); err != nil {
			panic(common.ErrorInvalidRequest("Create room request", err))
		}

		gameService := service.NewGameService(appCtx)
		room, err := gameService.CreateRoom(ginCtx.Request.Context(), createRoomReq)
		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse(room))
	}
}
