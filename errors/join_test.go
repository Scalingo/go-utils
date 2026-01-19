package errors

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Custom error types for testing
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

var ErrNotFound = errors.New("not found")

type testContextKey string

const (
	testKey1    testContextKey = "testKey1"
	testKey2    testContextKey = "testKey2"
	testKeyRoot testContextKey = "testKeyRoot"
)

func TestJoin_RootCtxFromFirstError(t *testing.T) {
	t.Run("It should get the root context from the first errctx error when joining 2 errctx errors", func(t *testing.T) {
		// Given
		ctx1 := context.WithValue(t.Context(), testKey1, "value1")
		ctx1 = context.WithValue(ctx1, testKeyRoot, "first")
		err1 := New(ctx1, "first error")

		ctx2 := context.WithValue(t.Context(), testKey2, "value2")
		ctx2 = context.WithValue(ctx2, testKeyRoot, "second")
		err2 := New(ctx2, "second error")

		// When
		joinedErr := Join(err1, err2)
		require.Error(t, joinedErr)

		// Then
		rootCtx := RootCtxOrFallback(t.Context(), joinedErr)
		assert.NotNil(t, rootCtx)
		assert.Equal(t, "value1", rootCtx.Value(testKey1))
		assert.Equal(t, "first", rootCtx.Value(testKeyRoot))
		// The second error's context should not be accessible
		assert.Nil(t, rootCtx.Value(testKey2))
	})

	t.Run("It should get the root context from nested errctx errors in the first error when joining", func(t *testing.T) {
		// Given
		// First error with nested context
		rootCtx1 := context.WithValue(t.Context(), testKeyRoot, "root_value")
		err1 := New(rootCtx1, "base error")

		wrapCtx1 := context.WithValue(t.Context(), testKey1, "wrap_value")
		err1 = Wrap(wrapCtx1, err1, "wrapped error")

		// Second error with different context
		ctx2 := context.WithValue(t.Context(), testKey2, "value2")
		err2 := New(ctx2, "second error")

		// When
		joinedErr := Join(err1, err2)
		require.Error(t, joinedErr)

		// Then - should get the root context from the first error
		rootCtx := RootCtxOrFallback(t.Context(), joinedErr)
		assert.NotNil(t, rootCtx)
		assert.Equal(t, "root_value", rootCtx.Value(testKeyRoot))
		// Should not have values from wrapping or second error
		assert.Nil(t, rootCtx.Value(testKey1))
		assert.Nil(t, rootCtx.Value(testKey2))
	})
}

func TestJoin_Is(t *testing.T) {
	t.Run("It should match error value with Is when wrapped in ErrCtx and joined", func(t *testing.T) {
		// Given
		ctx1 := context.WithValue(t.Context(), testKey1, "value1")
		err1 := Wrap(ctx1, ErrNotFound, "wrapped not found")

		ctx2 := context.WithValue(t.Context(), testKey2, "value2")
		err2 := New(ctx2, "second error")

		// When
		joinedErr := Join(err1, err2)
		require.Error(t, joinedErr)

		// Then
		assert.True(t, Is(joinedErr, ErrNotFound))
	})

	t.Run("It should match error value with Is when in second joined error", func(t *testing.T) {
		// Given
		ctx1 := context.WithValue(t.Context(), testKey1, "value1")
		err1 := New(ctx1, "first error")

		ctx2 := context.WithValue(t.Context(), testKey2, "value2")
		err2 := Wrap(ctx2, ErrNotFound, "wrapped not found in second")

		// When
		joinedErr := Join(err1, err2)
		require.Error(t, joinedErr)

		// Then
		assert.True(t, Is(joinedErr, ErrNotFound))
	})

	t.Run("It should match multiple error values with Is when both are wrapped in ErrCtx and joined", func(t *testing.T) {
		// Given
		ErrDatabase := errors.New("database error")
		ErrNetwork := errors.New("network error")

		ctx1 := context.WithValue(t.Context(), testKey1, "value1")
		err1 := Wrap(ctx1, ErrDatabase, "wrapped database error")

		ctx2 := context.WithValue(t.Context(), testKey2, "value2")
		err2 := Wrap(ctx2, ErrNetwork, "wrapped network error")

		// When
		joinedErr := Join(err1, err2)
		require.Error(t, joinedErr)

		// Then
		require.ErrorIs(t, joinedErr, ErrDatabase)
		require.ErrorIs(t, joinedErr, ErrNetwork)
	})
}

func TestJoin_As(t *testing.T) {
	t.Run("It should match error value with As when wrapped in ErrCtx and joined", func(t *testing.T) {
		// Given
		customErr := &CustomError{Code: 404, Message: "resource not found"}
		ctx1 := context.WithValue(t.Context(), testKey1, "value1")
		err1 := Wrap(ctx1, customErr, "wrapped custom error")

		ctx2 := context.WithValue(t.Context(), testKey2, "value2")
		err2 := New(ctx2, "second error")

		// When
		joinedErr := Join(err1, err2)
		require.Error(t, joinedErr)

		// Then
		var target *CustomError
		assert.True(t, As(joinedErr, &target))
		assert.NotNil(t, target)
		assert.Equal(t, 404, target.Code)
		assert.Equal(t, "resource not found", target.Message)
	})

	t.Run("It should match error value with As when in second joined error", func(t *testing.T) {
		// Given
		ctx1 := context.WithValue(t.Context(), testKey1, "value1")
		err1 := New(ctx1, "first error")

		customErr := &CustomError{Code: 500, Message: "internal error"}
		ctx2 := context.WithValue(t.Context(), testKey2, "value2")
		err2 := Wrap(ctx2, customErr, "wrapped custom error in second")

		// When
		joinedErr := Join(err1, err2)
		require.Error(t, joinedErr)

		// Then
		var target *CustomError
		assert.True(t, As(joinedErr, &target))
		assert.NotNil(t, target)
		assert.Equal(t, 500, target.Code)
		assert.Equal(t, "internal error", target.Message)
	})
}
