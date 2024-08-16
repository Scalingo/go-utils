package pagination

const defaultPageSize = 20

type PageRequest struct {
	page     int // page requested (default 1)
	pageSize int // Number of items per page
}

func NewPageRequest(page, pageSize int) PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	return PageRequest{
		page:     page,
		pageSize: pageSize,
	}
}

func (p PageRequest) Page() int {
	return p.page
}

func (p PageRequest) PageSize() int {
	return p.pageSize
}

func (p PageRequest) Limit() int32 {
	return int32(p.pageSize)
}

func (p PageRequest) Offset() int32 {
	return int32((p.page - 1) * p.pageSize)
}
