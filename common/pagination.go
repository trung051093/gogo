package common

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

type Pagination struct {
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
	Total  int64 `json:"total"`
	Offset int   `json:"offset"`
}

func (p *Pagination) Paginate() error {
	if p.Page <= 0 {
		p.Page = DefaultPage
	}

	if p.Limit <= 0 || p.Limit >= 1001 {
		p.Limit = DefaultLimit
	}

	p.Offset = (p.Page - 1) * p.Limit

	return nil
}
