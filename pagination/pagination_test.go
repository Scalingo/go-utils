package pagination

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type testTypeA []int
type testTypeB []string

func TestConvertData(t *testing.T) {
	t.Run("convert paginated data from type a to type b", func(t *testing.T) {
		// Given
		initialData := testTypeA{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		paginatedA := NewPaginated(initialData, NewPageRequest(1, 10), 10)

		// When
		paginatedB := ConvertData(paginatedA, func(data testTypeA) testTypeB {
			var convertedData testTypeB
			for _, v := range data {
				convertedData = append(convertedData, fmt.Sprintf("%v", v))
			}
			return convertedData
		})

		// Then
		expectedData := testTypeB{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
		require.Equal(t, Paginated[testTypeB]{
			Data: expectedData,
			Pagination: Pagination{
				CurrentPage: 1,
				PageSize:    10,
				PrevPage:    1,
				NextPage:    1,
				TotalPages:  1,
				TotalCount:  10,
			},
		}, paginatedB)
	})
}

func TestNewPaginated(t *testing.T) {
	t.Run("new paginated pageSize 10 empty", func(t *testing.T) {
		// Given
		data := testTypeA{}
		pageRequest := NewPageRequest(1, 10)
		totalCount := int64(0)

		// When
		got := NewPaginated(data, pageRequest, totalCount)

		// Then
		require.Equal(t, Paginated[testTypeA]{
			Data: testTypeA{},
			Pagination: Pagination{
				CurrentPage: 1,
				PageSize:    10,
				PrevPage:    1,
				NextPage:    1,
				TotalPages:  1,
				TotalCount:  0,
			},
		}, got)
	})

	t.Run("new paginated pageSize 10 one full page", func(t *testing.T) {
		// Given
		data := testTypeA{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		pageRequest := NewPageRequest(1, 10)
		totalCount := int64(10)

		// When
		got := NewPaginated(data, pageRequest, totalCount)

		// Then
		require.Equal(t, Paginated[testTypeA]{
			Data: testTypeA{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Pagination: Pagination{
				CurrentPage: 1,
				PageSize:    10,
				PrevPage:    1,
				NextPage:    1,
				TotalPages:  1,
				TotalCount:  10,
			},
		}, got)
	})
}
