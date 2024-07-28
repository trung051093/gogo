package dto

import "gogo/modules/tracker/entity"

type GameAddUserDto struct {
	GameId       string           `json:"game_id" uri:"game_id" binding:"required" validate:"required"`
	Users        []*UserCreateDto `json:"users"`
	DefaultBuyIn int              `json:"default_buy_in" validate:"required"`
}

func (dt *GameAddUserDto) Validate() error {
	return nil
}

func (dt *GameAddUserDto) ToUserEntities() []*entity.User {
	var users []*entity.User
	for _, user := range dt.Users {
		users = append(users, user.ToEntity())
	}
	return users
}
