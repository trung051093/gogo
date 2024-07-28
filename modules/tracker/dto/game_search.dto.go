package dto

import (
	"gogo/common"
)

type GameSearchDto struct {
	*common.CursorPagination
}

func (u *GameSearchDto) Validate() error {
	return nil
}

func (u *GameSearchDto) ToQuery() common.GormQuery {
	gormQuery := common.NewGormQuery()
	return gormQuery
}

func (u *GameSearchDto) ToPagination() *common.CursorPagination {
	return u.CursorPagination
}
