package pagination

import "math"

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"per_page"`
	PrevPage    int   `json:"prev_page"`
	NextPage    int   `json:"next_page"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

type Paginated[T interface{}] struct {
	Data T          `json:"data"`
	Meta Pagination `json:"meta"`
}

func NewPaginated[T interface{}](data T, pageRequest PageRequest, totalCount int64) Paginated[T] {
	prevPage := pageRequest.page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageRequest.pageSize)))
	if totalPages == 0 {
		// We always want at least one page, even if it is empty
		totalPages = 1
	}

	nextPage := pageRequest.page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	return Paginated[T]{
		Data: data,
		Meta: Pagination{
			CurrentPage: pageRequest.page,
			PageSize:    pageRequest.pageSize,
			PrevPage:    prevPage,
			NextPage:    nextPage,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
}
