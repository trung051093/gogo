package routes

import (
	"gogo/components/appctx"
	decorator "gogo/decorators"
	"gogo/middleware"
	"gogo/modules/auth"
	"gogo/modules/file"
	"gogo/modules/user"
	"time"

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
		v1.POST("/auth/logout", middleware.JWTRequireHandler(appCtx), auth.LogoutUserHandler(appCtx))
		v1.POST("/auth/forgot-password", auth.ForgotPasswordUserHandler(appCtx))
		v1.POST("/auth/reset-password", auth.ResetPasswordUserHandler(appCtx))

		// file
		v1.GET("/file/presign-url", file.GetUploadPresignedUrl(appCtx))
	}
}
