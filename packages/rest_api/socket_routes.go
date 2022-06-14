package main

import (
	"user_management/components/appctx"
	"user_management/middleware"

	"github.com/gin-gonic/gin"
)

func socketRoutes(appCtx appctx.AppContext, router *gin.Engine) {
	socketService := appCtx.GetSocketService()
	socketServer := socketService.GetServer()

	socketRoute := router.Group("/socket.io")
	{
		socketRoute.Use(middleware.CorsMiddleware("*"))
		socketRoute.GET("/*any", gin.WrapH(socketServer))
		socketRoute.POST("/*any", gin.WrapH(socketServer))
	}

}
