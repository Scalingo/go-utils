package pagination

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testTypeA []int

func TestNewPaginated(t *testing.T) {
	t.Run("new paginated perPage 10 empty", func(t *testing.T) {
		// Given
		data := testTypeA{}
		pageRequest := NewRequest(1, 10)
		totalCount := int64(0)

		// When
		got := New(data, pageRequest, totalCount)

		// Then
		require.Equal(t, Paginated[testTypeA]{
			Data: testTypeA{},
			Meta: Meta{
				CurrentPage: 1,
				PerPage:     10,
				PrevPage:    1,
				NextPage:    1,
				TotalPages:  1,
				TotalCount:  0,
			},
		}, got)
	})

	t.Run("new paginated perPage 10 one full page", func(t *testing.T) {
		// Given
		data := testTypeA{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		pageRequest := NewRequest(1, 10)
		totalCount := int64(10)

		// When
		got := New(data, pageRequest, totalCount)

		// Then
		require.Equal(t, Paginated[testTypeA]{
			Data: testTypeA{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Meta: Meta{
				CurrentPage: 1,
				PerPage:     10,
				PrevPage:    1,
				NextPage:    1,
				TotalPages:  1,
				TotalCount:  10,
			},
		}, got)
	})
}
