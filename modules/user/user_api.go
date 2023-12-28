package user

import (
	"gogo/common"
	"gogo/components/appctx"
	usermodel "gogo/modules/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary     Create user
// @Description create user
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       user body     usermodel.UserCreateDto      true "Add user"
// @Success     200  {object} common.Response{data=string} "desc"
// @Failure     400  {object} common.AppError
// @Router      /api/v1/user [post]
func CreateUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		createDto := &usermodel.UserCreateDto{}
		if err := common.ParseRequest[usermodel.UserCreateDto](appCtx, ginCtx)(createDto); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.User{}.EntityName(), err))
		}

		userService := NewUserServiceWithAppCtx(appCtx)
		userId, err := userService.Create(ginCtx.Request.Context(), createDto)
		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse(userId))
	}
}

// UpdateUser godoc
// @Summary     Update an user
// @Description update user
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       user body     usermodel.UserUpdateDto      true "Add account"
// @Success     200  {object} common.Response{data=string} "desc"
// @Failure     400  {object} common.AppError
// @Router      /api/v1/user/{id} [patch]
func UpdateUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		updateDto := &usermodel.UserUpdateDto{}
		userId := ginCtx.Param("id")
		if err := common.ParseRequest[usermodel.UserUpdateDto](appCtx, ginCtx)(updateDto); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.User{}.EntityName(), err))
		}

		userService := NewUserServiceWithAppCtx(appCtx)
		if _, err := userService.UpdateByID(ginCtx.Request.Context(), userId, updateDto); err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse("OK"))
	}
}

// GetUser godoc
// @Summary     Get an user
// @Description get user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     int                                  true "User ID"
// @Success     200 {object} common.Response{data=usermodel.User} "desc"
// @Failure     400 {object} common.AppError
// @Router      /api/v1/user/{id} [get]
func GetUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		userId := ginCtx.Param("id")
		userService := NewUserServiceWithAppCtx(appCtx)
		user, err := userService.FindById(ginCtx.Request.Context(), userId)
		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(user))
	}
}

// GetListUser godoc
// @Summary     Get list of user
// @Description get string by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       fields    query    int                                                                                 false "fields"
// @Param       before    query    string                                                                              false "before cursor"
// @Param       after     query    string                                                                              false "after cursor"
// @Param       limit     query    int                                                                                 true  "limit"
// @Param       sortField query    string                                                                              false "sort by field"
// @Param       sortName  query    string                                                                              false "sort by field"
// @Success     200       {object} common.ResponsePagination{data=[]usermodel.User,pagination=common.CursorPagination} "desc"
// @Failure     400       {object} common.AppError
// @Router      /api/v1/users [get]
func ListUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		searchDto := &usermodel.UserSearchDto{
			CursorPagination: &common.CursorPagination{},
		}
		if err := common.ParseRequest[usermodel.UserSearchDto](appCtx, ginCtx)(searchDto); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.User{}.EntityName(), err))
		}

		userService := NewUserServiceWithAppCtx(appCtx)
		data, filter, paging, err := userService.SearchCursorPaging(ginCtx.Request.Context(), searchDto)
		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponsePagination(data, paging, filter))
	}
}

// DeleteUser godoc
// @Summary     Delete an user
// @Description delete user
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id  path     int                        true "User ID"
// @Success     200 {object} common.Response{data=bool} "desc"
// @Failure     400 {object} common.AppError
// @Router      /api/v1/user/{id} [delete]
func DeleteUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		userId := ginCtx.Param("id")
		userService := NewUserServiceWithAppCtx(appCtx)

		if err := userService.DeleteById(ginCtx.Request.Context(), userId); err != nil {
			panic(common.ErrorCannotDeleteEntity(usermodel.User{}.EntityName(), err))
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(true))
	}
}

// SearchUser godoc
// @Summary     Search an user
// @Description search user
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       query     query    string                                             true  "query"
// @Param       lastIndex query    string                                             false "lastIndex"
// @Param       sortField query    string                                             false "sort by field"
// @Param       sortName  query    string                                             false "sort by field"
// @Param       id        path     int                                                false "User ID"
// @Success     200       {object} common.Response{data=usermodel.UserEsSearchResult} "desc"
// @Failure     400       {object} common.AppError
// @Router      /api/v1/user/search [get]
func SearchUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		searchDto := &usermodel.UserEsSearchDto{}
		if err := common.ParseRequest[usermodel.UserEsSearchDto](appCtx, ginCtx)(searchDto); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.User{}.EntityName(), err))
		}

		userService := NewUserServiceWithAppCtx(appCtx)
		data, err := userService.EsSearch(ginCtx.Request.Context(), searchDto)

		if err != nil {
			panic(common.ErrorCannotFoundEntity(usermodel.User{}.EntityName(), err))
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}
