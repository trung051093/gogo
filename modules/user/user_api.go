package user

import (
	"net/http"
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

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		userService := NewUserService(userRepo)

		if err := userService.CreateUser(&newData); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{"data": true})
	}
}
