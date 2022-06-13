package main

import (
	"user_management/components/appctx"

	"github.com/gin-gonic/gin"
)

func socketRoutes(appCtx appctx.AppContext, router *gin.Engine) {
	socketService := appCtx.GetSocketService()
	socketServer := socketService.GetServer()

	router.GET("/socket.io/*any", func(ginCtx *gin.Context) {
		gin.WrapH(socketServer)
	})
	// Method 2 using server.ServerHTTP(Writer, Request) and also you can simply this by using gin.WrapH
	router.POST("/socket.io/*any", func(ginCtx *gin.Context) {
		gin.WrapH(socketServer)
	})
}
