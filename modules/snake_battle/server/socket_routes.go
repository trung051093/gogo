package server

import (
	"gogo/modules/snake_battle/middleware"

	"github.com/gin-gonic/gin"
)

func (s *SnakeBattleServer) createSocketRoutes() {
	socketService := s.appCtx.GetSocketService()
	socketServer := socketService.GetServer()

	socketRoute := s.ginEngine.Group("/socket.io")
	{
		socketRoute.Use(middleware.CorsMiddleware("*"))
		socketRoute.GET("/*any", gin.WrapH(socketServer))
		socketRoute.POST("/*any", gin.WrapH(socketServer))
	}
}
