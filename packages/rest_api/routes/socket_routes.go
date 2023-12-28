package routes

import (
	"gogo/components/appctx"
	"gogo/middleware"

	"github.com/gin-gonic/gin"
)

func SocketRoutes(appCtx appctx.AppContext, router *gin.Engine) {
	socketService := appCtx.GetSocketService()
	socketServer := socketService.GetServer()

	socketRoute := router.Group("/socket.io")
	{
		socketRoute.Use(middleware.CorsMiddleware("*"))
		socketRoute.GET("/*any", gin.WrapH(socketServer))
		socketRoute.POST("/*any", gin.WrapH(socketServer))
	}

}
