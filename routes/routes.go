package routes

import (
	"time"
	"user_management/components/appctx"
	decorator "user_management/decorators"
	"user_management/modules/auth"
	"user_management/modules/file"
	"user_management/modules/user"

	"github.com/gin-gonic/gin"
)

func MainRoutes(appCtx appctx.AppContext, router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// user
		v1.POST("/user", user.CreateUserHandler(appCtx))
		v1.PATCH("/user/:id", user.UpdateUserHandler(appCtx))
		v1.DELETE("/user/:id", user.DeleteUserHandler(appCtx))
		v1.GET("/user/:id", user.GetUserHandler(appCtx))
		v1.GET("/users", decorator.CacheRequest(appCtx, "user", 1*time.Minute)(user.ListUserHandler))
		v1.GET("/user/search", user.SearchUserHandler(appCtx))

		// authentication
		v1.POST("/auth/register", auth.RegisterUserHandler(appCtx))
		v1.POST("/auth/login", auth.LoginUserHandler(appCtx))

		// file
		v1.GET("/file/presign-url", file.GetUploadPresignedUrl(appCtx))
	}
}
