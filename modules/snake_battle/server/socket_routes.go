package server

import (
	"context"
	"gogo/modules/snake_battle/service"

	"github.com/gin-gonic/gin"
)

func (s *SnakeBattleServer) createSocketRoutes() {
	socketSvc := s.appCtx.GetSocketService()
	gameSvc := service.NewGameService(s.appCtx)
	gameSvc.GameSocketListener(context.Background())
	s.ginEngine.GET("/socket.io/*any", gin.WrapH(socketSvc.GetServer()))
	s.ginEngine.POST("/socket.io/*any", gin.WrapH(socketSvc.GetServer()))
}
