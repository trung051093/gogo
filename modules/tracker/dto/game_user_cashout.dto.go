package dto

type GameUserCashoutDto struct {
	GameId  string `json:"game_id" uri:"game_id" binding:"required" validate:"required"`
	UserId  string `json:"user_id" uri:"user_id" binding:"required"`
	Cashout int    `json:"cashout"`
}

func (dt *GameUserCashoutDto) Validate() error {
	return nil
}
