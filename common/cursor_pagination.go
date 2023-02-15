package common

import (
	"strings"

	"github.com/pilagod/gorm-cursor-paginator/v2/cursor"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

// =====cursor pagination=====
// default is cursor pagination
type CursorPagination struct {
	Paginator *paginator.Paginator `json:"-"`
	Cursor    paginator.Cursor     `json:"cursor"`
	Limit     int                  `json:"limit"`
	Total     int64                `json:"total"`
	Order     string               `json:"-"`
	Rules     []paginator.Rule     `json:"-"`
}

func (p *CursorPagination) Paginate() {
	p.Rules = make([]paginator.Rule, 0)
	if p.Limit <= 0 || p.Limit >= 100 {
		p.Limit = DefaultLimit
	}
}

func NewCursorPagination(config *paginator.Config, configRules map[string]paginator.Rule) *CursorPagination {
	p := &CursorPagination{
		Limit: DefaultLimit,
	}

	if config == nil {
		config = &paginator.Config{
			Rules: []paginator.Rule{
				{
					Key:             "CreatedAt",
					Order:           paginator.DESC,
					NULLReplacement: "1970-01-01",
				},
			},
			Limit: 10,
			Order: paginator.DESC,
		}
	}
	paginator := p.NewPaginator(config, configRules)
	p.Paginator = paginator
	return p
}

func (p *CursorPagination) BuildOrderRule(mapKey map[string]paginator.Rule) {
	if p.Order != "" {
		listOrders := strings.Split(p.Order, ",")
		for _, orderKey := range listOrders {
			rule, ok := mapKey[orderKey]
			if ok {
				p.Rules = append(p.Rules, rule)
			}
		}
	}
}

func (p *CursorPagination) NewPaginator(config *paginator.Config, configRules map[string]paginator.Rule) *paginator.Paginator {
	// set default setting and rules
	p.Paginate()
	if configRules != nil {
		p.BuildOrderRule(configRules)
	}

	opts := []paginator.Option{config}
	if len(p.Rules) > 0 {
		opts = append(opts, paginator.WithRules(p.Rules...))
	}
	if p.Limit != 0 {
		opts = append(opts, paginator.WithLimit(p.Limit))
	}
	if p.Cursor.After != nil {
		opts = append(opts, paginator.WithAfter(*p.Cursor.After))
	}
	if p.Cursor.Before != nil {
		opts = append(opts, paginator.WithBefore(*p.Cursor.Before))
	}
	p.Paginator = paginator.New(opts...)
	return p.Paginator
}

func (p *CursorPagination) GetPaginator() *paginator.Paginator {
	if p.Paginator == nil {
		p.Paginator = p.NewPaginator(&paginator.Config{
			Keys:  []string{"ID"},
			Limit: 10,
			// Order here will apply to keys without order specified.
			Order: paginator.DESC,
		}, nil)
	}
	return p.Paginator
}

func (p *CursorPagination) GetLimit() int {
	return p.Limit
}

func (p *CursorPagination) SetLimit(limit int) {
	p.Limit = limit
}

func (p *CursorPagination) GetTotal() int64 {
	return p.Total
}

func (p *CursorPagination) SetTotal(total int64) {
	p.Total = total
}

func (p *CursorPagination) GetOffset() int {
	return 0
}

func (p *CursorPagination) SetOffset(offset int) {
}

func (p *CursorPagination) SetCursor(c cursor.Cursor) {
	p.Cursor = c
}

func (p *CursorPagination) GetCursor() cursor.Cursor {
	return p.Cursor
}
