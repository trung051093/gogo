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
		v1.POST("/user", middleware.JWTRequireHandler(appCtx), user.CreateUserHandler(appCtx))
		v1.PATCH("/user/:id", middleware.JWTRequireHandler(appCtx), user.UpdateUserHandler(appCtx))
		v1.DELETE("/user/:id", middleware.JWTRequireHandler(appCtx), user.DeleteUserHandler(appCtx))
		v1.GET("/user/:id", middleware.JWTRequireHandler(appCtx), user.GetUserHandler(appCtx))
		v1.GET("/users", middleware.JWTRequireHandler(appCtx), user.ListUserHandler(appCtx))
		v1.GET("/users-cache", middleware.JWTRequireHandler(appCtx), decorator.CacheRequest(appCtx, "user", time.Minute)(user.ListUserHandler))
		v1.GET("/user/search", middleware.JWTRequireHandler(appCtx), user.SearchUserHandler(appCtx))

		// authentication
		v1.POST("/auth/register", auth.RegisterUserHandler(appCtx))
		v1.POST("/auth/login", auth.LoginUserHandler(appCtx))
		v1.POST("/auth/logout", middleware.JWTRequireHandler(appCtx), auth.LogoutUserHandler(appCtx))
		v1.POST("/auth/forgot-password", auth.ForgotPasswordUserHandler(appCtx))
		v1.POST("/auth/reset-password", auth.ResetPasswordUserHandler(appCtx))

		v1.GET("/auth/google/login", auth.GoogleLoginUserHandler(appCtx))
		v1.GET("/auth/google/callback", auth.GoogleCallbackUserHandler(appCtx))

		// file
		v1.GET("/file/presign-url", file.GetUploadPresignedUrl(appCtx))
	}
}
