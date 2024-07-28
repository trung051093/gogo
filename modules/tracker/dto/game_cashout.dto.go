package dto

type UserCashoutDto struct {
	UserId  string `json:"user_id"`
	Cashout int    `json:"cashout"`
}

type GameCashoutDto struct {
	GameId         string            `json:"game_id" uri:"game_id" binding:"required" validate:"required"`
	CashoutSummary []*UserCashoutDto `json:"cashout_summary"`
}

func (dt *GameCashoutDto) Validate() error {
	return nil
}
