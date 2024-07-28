package dto

import "gogo/modules/tracker/entity"

type GameCreateDto struct {
	Name         string           `json:"name" validate:"required"`
	Users        []*UserCreateDto `json:"users"`
	DefaultBuyIn int              `json:"default_buy_in" validate:"required"`
}

func (dt *GameCreateDto) Validate() error {
	return nil
}

func (dt *GameCreateDto) ToEntity() *entity.Game {
	return &entity.Game{
		Name: dt.Name,
	}
}

func (dt *GameCreateDto) ToUserEntities() []*entity.User {
	var users []*entity.User
	for _, user := range dt.Users {
		users = append(users, user.ToEntity())
	}
	return users
}
