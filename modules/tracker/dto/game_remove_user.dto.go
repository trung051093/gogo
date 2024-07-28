package dto

type GameRemoveUserDto struct {
	GameId string `json:"game_id" uri:"game_id" binding:"required" validate:"required"`
	UserId string `json:"user_id" uri:"user_id" binding:"required" validate:"required"`
}

func (dt *GameRemoveUserDto) Validate() error {
	return nil
}
