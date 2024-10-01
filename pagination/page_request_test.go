package pagination

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPageRequest(t *testing.T) {
	t.Run("creates new page request correctly initialised", func(t *testing.T) {
		got := NewPageRequest(2, 22)
		require.Equal(t, PageRequest{
			Page:     2,
			PageSize: 22,
		}, got)
	})

	t.Run("creates new page request with defaults when badly initialised", func(t *testing.T) {
		got := NewPageRequest(0, 0)
		require.Equal(t, PageRequest{
			Page:     1,
			PageSize: 20,
		}, got)
	})
}

func TestPageRequest_Limit(t *testing.T) {
	t.Run("get limit same size as initialised pageSize", func(t *testing.T) {
		// Given
		p := NewPageRequest(1, 20)

		// When
		limit := p.Limit()

		// Then
		require.Equal(t, 20, limit)
	})
}

func TestPageRequest_Offset(t *testing.T) {
	t.Run("get offset 0 if page is 1", func(t *testing.T) {
		// Given
		p := NewPageRequest(1, 20)

		// When
		offset := p.Offset()

		// Then
		require.Equal(t, 0, offset)
	})
	t.Run("get offset 20 if page is 2", func(t *testing.T) {
		// Given
		p := NewPageRequest(2, 20)

		// When
		offset := p.Offset()

		// Then
		require.Equal(t, 20, offset)
	})
	t.Run("get offset 40 if page is 3", func(t *testing.T) {
		// Given
		p := NewPageRequest(3, 20)

		// When
		offset := p.Offset()

		// Then
		require.Equal(t, 40, offset)
	})
	t.Run("get offset 30 if page is 2 and page size is 30", func(t *testing.T) {
		// Given
		p := NewPageRequest(2, 30)

		// When
		offset := p.Offset()

		// Then
		require.Equal(t, 30, offset)
	})
}
