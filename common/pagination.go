package common

type Pagination struct {
	Page   int
	Limit  int
	Total  int64
	Offset int
}

func (p *Pagination) Paginate() error {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit >= 1001 {
		p.Limit = 10
	}

	p.Offset = (p.Page - 1) * p.Limit

	return nil
}
