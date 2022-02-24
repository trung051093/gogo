package user

import (
	"net/http"
	"strconv"
	"user_management/common"
	component "user_management/components"
	usermodel "user_management/modules/user/model"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ctx *gin.Context
}

func CreateUserHandler(appCtx component.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var newData usermodel.UserCreate

		if err := ginCtx.ShouldBind(&newData); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := appCtx.GetValidator()
		if err := validate.Struct(&newData); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		if err := userService.CreateUser(&newData); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{"data": true})
	}
}

func UpdateUserHandler(appCtx component.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var updateData usermodel.UserUpdate

		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ginCtx.ShouldBind(&updateData); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := appCtx.GetValidator()
		if err := validate.Struct(&updateData); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		if err := userService.UpdateUser(uint(id), &updateData); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{"data": true})
	}
}

func GetUserHandler(appCtx component.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		user, err := userService.GetUser(uint(id))

		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func ListUserHandler(appCtx component.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var paging common.Pagination

		if err := ginCtx.ShouldBind(&paging); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := paging.Paginate(); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		data, err := userService.SearchUsers(map[string]interface{}{}, &paging)

		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{"data": data, "pagination": paging})
	}
}

func DeleteUserHandler(appCtx component.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		if err := userService.DeleteUser(uint(id)); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{"data": true})
	}
}
