package server

import (
	"gogo/modules/tracker/api"
	"gogo/modules/tracker/middleware"
)

func (s *TrackerServer) createMainRoutes() {

	v1 := s.ginEngine.Group("/api/v1")
	{
		// user
		v1.GET("/users", middleware.Authenticaticator(s.appCtx), api.ListUserHandler(s.appCtx))

		// authentication
		v1.POST("/auth/register", api.RegisterUserHandler(s.appCtx))
		v1.POST("/auth/login", api.LoginUserHandler(s.appCtx))
		v1.POST("/auth/logout", middleware.Authenticaticator(s.appCtx), api.LogoutUserHandler(s.appCtx))

		// game
		v1.GET("/games", middleware.Authenticaticator(s.appCtx), api.ListGameHandler(s.appCtx))
		v1.POST("/games", middleware.Authenticaticator(s.appCtx), middleware.AdminRequired(s.appCtx), api.NewGameHandler(s.appCtx))
		v1.GET("/games/:game_id", middleware.Authenticaticator(s.appCtx), middleware.AdminRequired(s.appCtx), api.GetGameByIDHandler(s.appCtx))
		v1.POST("/games/:game_id/users", middleware.Authenticaticator(s.appCtx), middleware.AdminRequired(s.appCtx), api.AddUserGameHandler(s.appCtx))
		v1.POST("/games/:game_id/summary", middleware.Authenticaticator(s.appCtx), middleware.AdminRequired(s.appCtx), api.SummaryHandler(s.appCtx))
		v1.DELETE("/games/:game_id/users/:user_id", middleware.Authenticaticator(s.appCtx), middleware.AdminRequired(s.appCtx), api.RemoveUserHandler(s.appCtx))
		v1.POST("/games/:game_id/users/:user_id/buyin", middleware.Authenticaticator(s.appCtx), middleware.AdminRequired(s.appCtx), api.BuyinGameHandler(s.appCtx))
		v1.POST("/games/:game_id/users/:user_id/cashout", middleware.Authenticaticator(s.appCtx), middleware.AdminRequired(s.appCtx), api.CashoutGameHandler(s.appCtx))
	}
}
