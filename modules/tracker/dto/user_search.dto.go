package dto

import (
	"gogo/common"
)

type UserSearchDto struct {
	*common.CursorPagination
}

func (u *UserSearchDto) Validate() error {
	return nil
}

func (u *UserSearchDto) ToQuery() common.GormQuery {
	gormQuery := common.NewGormQuery()
	return gormQuery
}

func (u *UserSearchDto) ToPagination() *common.CursorPagination {
	return u.CursorPagination
}
