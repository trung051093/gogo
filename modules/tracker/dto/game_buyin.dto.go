package dto

type GameBuyInDto struct {
	GameId string `json:"game_id" uri:"game_id" binding:"required" validate:"required"`
	UserId string `json:"user_id" uri:"user_id" binding:"required" validate:"required"`
	BuyIn  int    `json:"buy_in" validate:"required"`
}

func (dt *GameBuyInDto) Validate() error {
	return nil
}
