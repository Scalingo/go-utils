package pagination

const defaultPageSize = 20

type PageRequest struct {
	Page     int // page requested (default 1)
	PageSize int // Number of items per page
}

func NewPageRequest(page, pageSize int) PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	return PageRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p PageRequest) Limit() int {
	return p.PageSize
}

func (p PageRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}
