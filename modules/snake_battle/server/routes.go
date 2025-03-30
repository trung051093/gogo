package server

import "gogo/modules/snake_battle/api"

func (s *SnakeBattleServer) createMainRoutes() {
	v1 := s.ginEngine.Group("/api/v1")
	{
		v1.POST("/rooms", api.CreateRoom(s.appCtx))
	}
}
