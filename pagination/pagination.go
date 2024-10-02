package pagination

import "math"

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"per_page"`
	PrevPage    int   `json:"prev_page"`
	NextPage    int   `json:"next_page"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"` // int64 to support result sets with >2.17B rows
}

type Paginated[T any] struct {
	Data T          `json:"data"`
	Meta Pagination `json:"meta"`
}

func New[T any](data T, pageRequest PageRequest, totalCount int64) Paginated[T] {
	prevPage := pageRequest.Page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageRequest.PageSize)))
	if totalPages == 0 {
		// We always want at least one page, even if it is empty
		totalPages = 1
	}

	nextPage := pageRequest.Page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	return Paginated[T]{
		Data: data,
		Meta: Pagination{
			CurrentPage: pageRequest.Page,
			PageSize:    pageRequest.PageSize,
			PrevPage:    prevPage,
			NextPage:    nextPage,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
}
