package user

import (
	"net/http"
	"strconv"
	"strings"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/hasher"
	usermodel "user_management/modules/user/model"

	"github.com/gin-gonic/gin"
)

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

		appConfig := appCtx.GetConfig()

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)
		hashService := hasher.NewHashService()

		passwordSalt := hashService.GenerateRandomString(appConfig.JWT.PasswordSaltLength)
		hashPassword := hashService.GenerateSHA256(newData.Password, passwordSalt)
		newData.Password = hashPassword

		userId, err := userService.CreateUserTrace(ginCtx.Request.Context(), &newData)

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
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)

		user, err := userService.GetUserTrace(ginCtx.Request.Context(), id)

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

		data, err := userService.SearchUsersTrace(ginCtx.Request.Context(), map[string]interface{}{}, &filter, &paging)

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
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)

		if _, err := userService.DeleteUser(ginCtx.Request.Context(), id); err != nil {
			panic(common.ErrorCannotDeleteEntity(usermodel.EntityName, err))
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func SearchUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var filter usermodel.UserFilter
		var paging common.Pagination

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

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
