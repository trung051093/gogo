package usermodel

import (
	"gogo/common"
)

type UserSearchDto struct {
	*common.CursorPagination
	Fields []string `json:"fields" query:"fields"`
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
