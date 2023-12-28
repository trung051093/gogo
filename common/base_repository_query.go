package common

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormQuery interface {
	GetQueryArgs(query string) []interface{}
	GetQueries() map[string][]interface{}
	GetJoins() []string
	GetGroupBy() string
	GetOrders() []string
	GetFieldSelect() string

	SetJoins(joinQuery string) GormQuery
	SetQuery(query string, agrs ...interface{}) GormQuery
	SetGroupBy(groupBy string) GormQuery
	SetOrder(field string, direction OrderDirection) GormQuery
	SetFieldSelect(selectFields string) GormQuery
}

type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
)

type gormQuery struct {
	cond         map[string][]interface{}
	join         []string
	groupBy      string
	orders       []string
	selectFields string
}

func NewGormQuery() GormQuery {
	return &gormQuery{
		cond: make(map[string][]interface{}),
	}
}

func (q *gormQuery) GetOrders() []string {
	return q.orders
}

func (q *gormQuery) SetOrder(field string, direction OrderDirection) GormQuery {
	q.orders = append(q.orders, fmt.Sprintf("%s %s", field, direction))
	return q
}

func (q *gormQuery) GetFieldSelect() string {
	return q.selectFields
}

func (q *gormQuery) SetFieldSelect(selectFields string) GormQuery {
	q.selectFields = selectFields
	return q
}

func (q *gormQuery) SetGroupBy(groupBy string) GormQuery {
	q.groupBy = groupBy
	return q
}

func (q *gormQuery) SetJoins(joinQuery string) GormQuery {
	q.join = append(q.join, joinQuery)
	return q
}

func (q *gormQuery) SetQuery(query string, agrs ...interface{}) GormQuery {
	q.cond[query] = agrs
	return q
}

func (q *gormQuery) GetQueryArgs(query string) []interface{} {
	return q.cond[query]
}

func (q *gormQuery) GetQueries() map[string][]interface{} {
	return q.cond
}

func (q *gormQuery) GetJoins() []string {
	return q.join
}

func (q *gormQuery) GetGroupBy() string {
	return q.groupBy
}

func (r *repository[E]) HandleQuery(stmt *gorm.DB, q GormQuery) *gorm.DB {
	// set join query
	for _, joinQuery := range q.GetJoins() {
		stmt = stmt.Joins(joinQuery)
	}

	// build condition
	for query, args := range q.GetQueries() {
		if conds := stmt.Statement.BuildCondition(query, args...); len(conds) > 0 {
			stmt.Statement.AddClause(clause.Where{Exprs: conds})
		}
	}

	// set group by
	if q.GetGroupBy() != "" {
		stmt = stmt.Group(q.GetGroupBy())
	}

	// set order
	if len(q.GetOrders()) >= 0 {
		for _, order := range q.GetOrders() {
			stmt = stmt.Order(order)
		}
	}

	// set field select
	if q.GetFieldSelect() != "" {
		stmt = stmt.Select(q.GetFieldSelect())
	}

	return stmt
}
