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

// GetListUser godoc
//	@Summary		Get list of user
//	@Description	get string by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			before	query		string																				false	"before cursor"
//	@Param			after	query		string																				false	"after cursor"
//	@Param			limit	query		int																					true	"limit"
//	@Success		200		{object}	common.ResponsePagination{data=[]entity.User,pagination=common.CursorPagination}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/users [get]
func ListUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		searchDto := &dto.UserSearchDto{
			CursorPagination: common.NewCursorPagination(nil, nil),
		}
		err := common.ParseRequest[dto.UserSearchDto](appCtx, ginCtx)(searchDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.User{}.EntityName(), err))

		userService := service.NewUserServiceWithAppCtx(appCtx)
		data, filter, paging, err := userService.SearchCursorPaging(ginCtx.Request.Context(), searchDto)
		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponsePagination(data, paging, filter))
	}
}
