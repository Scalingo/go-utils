package pagination

const defaultPageSize = 20

type Request struct {
	Page     int // page requested (default 1)
	PageSize int // Number of items per page
}

func NewRequest(page, pageSize int) Request {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	return Request{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Request) Limit() int {
	return p.PageSize
}

func (p Request) Offset() int {
	return (p.Page - 1) * p.PageSize
}
