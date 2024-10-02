package pagination

const defaultPerPage = 20

type Request struct {
	Page    int // page requested (default 1)
	PerPage int // Number of items per page
}

func NewRequest(page, perPage int) Request {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = defaultPerPage
	}
	return Request{
		Page:    page,
		PerPage: perPage,
	}
}

func (p Request) Limit() int {
	return p.PerPage
}

func (p Request) Offset() int {
	return (p.Page - 1) * p.PerPage
}
