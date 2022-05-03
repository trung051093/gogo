package user

import (
	"net/http"
	"strconv"
	"strings"
	"user_management/common"
	"user_management/components/appctx"
	usermodel "user_management/modules/user/model"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ctx *gin.Context
}

func CreateUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var newData usermodel.UserCreate

		if err := ginCtx.ShouldBind(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		validate := appCtx.GetValidator()
		if err := validate.Struct(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		userId, err := userService.CreateUser(ginCtx.Request.Context(), &newData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(userId))
	}
}

func UpdateUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var updateData usermodel.UserUpdate

		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		if err := ginCtx.ShouldBind(&updateData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		validate := appCtx.GetValidator()
		if err := validate.Struct(&updateData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		if err := userService.UpdateUser(ginCtx.Request.Context(), uint(id), &updateData); err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func GetUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		user, err := userService.GetUser(ginCtx.Request.Context(), uint(id))

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
	}
}

func ListUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var filter usermodel.UserFilter
		var paging common.Pagination

		filter.Order = ginCtx.Query("order")

		if err := filter.Process(); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		if fieldQuery := ginCtx.Query("fields"); fieldQuery != "" {
			filter.Fields = strings.Split(fieldQuery, ",")
		}

		if pageQuery, err := strconv.Atoi(ginCtx.Query("page")); err != nil {
			paging.Page = common.DefaultPage
		} else {
			paging.Page = pageQuery
		}

		if limitQuery, err := strconv.Atoi(ginCtx.Query("limit")); err != nil {
			paging.Limit = common.DefaultLimit
		} else {
			paging.Limit = limitQuery
		}

		if err := paging.Paginate(); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		data, err := userService.SearchUsers(ginCtx.Request.Context(), map[string]interface{}{}, &filter, &paging)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(data, paging, nil))
	}
}

func DeleteUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		if err := userService.DeleteUser(ginCtx.Request.Context(), uint(id)); err != nil {
			panic(common.ErrorCannotDeleteEntity(usermodel.EntityName, err))
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
