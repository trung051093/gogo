package common

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

type PagePagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
	Offset    int   `json:"offset"`
}

func NewPagePagination() *PagePagination {
	p := &PagePagination{
		Page:  DefaultPage,
		Limit: DefaultLimit,
	}
	p.CalculateTotalPage()
	p.Paginate()
	return p
}

func (p *PagePagination) Paginate() error {
	if p.Page <= 0 {
		p.Page = DefaultPage
	}

	if p.Limit <= 0 || p.Limit >= 1001 {
		p.Limit = DefaultLimit
	}

	p.Offset = (p.Page - 1) * p.Limit

	return nil
}

func roundU(val float64) int {
	if val > 0 {
		return int(val + 1.0)
	}
	return int(val)
}

func (p *PagePagination) CalculateTotalPage() {
	p.TotalPage = roundU(float64(p.Total) / float64(p.Limit))
}

func (p *PagePagination) GetTotal() int64 {
	return p.Total
}

func (p *PagePagination) SetTotal(total int64) {
	p.Total = total
	p.CalculateTotalPage()
}

func (p *PagePagination) GetLimit() int {
	return p.Limit
}

func (p *PagePagination) SetLimit(limit int) {
	p.Limit = limit
}

func (p *PagePagination) GetOffset() int {
	return p.Offset
}

func (p *PagePagination) SetOffset(offset int) {
	p.Offset = offset
}
