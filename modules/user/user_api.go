package user

import (
	"gogo/common"
	"gogo/components/appctx"
	"gogo/components/hasher"
	usermodel "gogo/modules/user/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary      Create user
// @Description  create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      usermodel.UserCreate       true  "Add user"
// @Success      200   {object}  common.Response{data=int}  "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/user [post]
func CreateUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var newData usermodel.UserCreate

		if err := ginCtx.ShouldBind(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		validator := appCtx.GetValidator()
		if err := validator.Struct(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		appConfig := appCtx.GetConfig()

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)
		hashService := hasher.NewHashService()

		passwordSalt := hashService.GenerateRandomString(appConfig.JWT.PasswordSaltLength)
		hashPassword := hashService.GenerateSHA256(newData.Password, passwordSalt)
		newData.Password = hashPassword
		newData.PasswordSalt = passwordSalt

		userId, err := userService.CreateUserTrace(ginCtx.Request.Context(), &newData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(userId))
	}
}

// UpdateUser godoc
// @Summary      Update an user
// @Description  update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      usermodel.UserUpdate        true  "Add account"
// @Success      200   {object}  common.Response{data=bool}  "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/user/{id} [patch]
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

		validator := appCtx.GetValidator()
		if err := validator.Struct(&updateData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		appConfig := appCtx.GetConfig()

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)
		hashService := hasher.NewHashService()

		passwordSalt := hashService.GenerateRandomString(appConfig.JWT.PasswordSaltLength)
		hashPassword := hashService.GenerateSHA256(updateData.Password, passwordSalt)
		updateData.Password = hashPassword

		if _, err := userService.UpdateUserTrace(ginCtx.Request.Context(), id, &updateData); err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(true))
	}
}

// GetUser godoc
// @Summary      Get an user
// @Description  get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int                                   true  "User ID"
// @Success      200  {object}  common.Response{data=usermodel.User}  "desc"
// @Failure      400        {object}  common.AppError
// @Router       /api/v1/user/{id} [get]
func GetUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)

		user, err := userService.GetUserTrace(ginCtx.Request.Context(), id)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(user))
	}
}

// GetListUser godoc
// @Summary      Get list of user
// @Description  get string by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        fields     query     int                                                                                                        false  "fields"
// @Param        page       query     int                                                                                                        true   "page"
// @Param        limit      query     int                                                                                                        true   "limit"
// @Param        sortField  query     string                                                                                                     false  "sort by field"
// @Param        sortName   query     string                                                                                                     false  "sort by field"
// @Success      200        {object}  common.ResponsePagination{data=[]usermodel.User,pagination=common.PagePagination,filter=usermodel.UserFilter}  "desc"
// @Failure      400  {object}  common.AppError
// @Router       /api/v1/users [get]
func ListUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		filter := &usermodel.UserFilter{}
		paging := &common.PagePagination{}

		filter.SortField = ginCtx.Query("sortField")
		filter.SortName = ginCtx.Query("sortName")

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
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)

		data, err := userService.SearchUsersTrace(ginCtx.Request.Context(), nil, filter, paging)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponsePagination(data, paging, filter))
	}
}

// DeleteUser godoc
// @Summary      Delete an user
// @Description  delete user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int                         true  "User ID"
// @Success      200  {object}  common.Response{data=bool}  "desc"
// @Failure      400  {object}  common.AppError
// @Router       /api/v1/user/{id} [delete]
func DeleteUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)

		if _, err := userService.DeleteUser(ginCtx.Request.Context(), id); err != nil {
			panic(common.ErrorCannotDeleteEntity(usermodel.EntityName, err))
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(true))
	}
}

// SearchUser godoc
// @Summary      Search an user
// @Description  search user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        query      query     int                                                 true   "query"
// @Param        lastIndex  query     int                                                 false  "lastIndex"
// @Param        sortField  query     string                                              false  "sort by field"
// @Param        sortName   query     string                                              false  "sort by field"
// @Param        id         path      int                                                 true   "User ID"
// @Success      200        {object}  common.Response{data=usermodel.UserEsSearchResult}  "desc"
// @Failure      400        {object}  common.AppError
// @Router       /api/v1/user/search [get]
func SearchUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var filter usermodel.UserFilter
		var paging common.PagePagination

		query := ginCtx.Query("query")
		lastIndex := ginCtx.Query("lastIndex")
		filter.SortField = ginCtx.Query("sortField")
		filter.SortName = ginCtx.Query("sortName")

		if err := filter.Process(); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
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
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)
		data, err := userService.EsSearchTrace(ginCtx.Request.Context(), query, lastIndex, &filter, &paging)

		if err != nil {
			panic(common.ErrorCannotFoundEntity(usermodel.EntityName, err))
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}
