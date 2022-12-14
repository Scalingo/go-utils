package errors

import (
	"context"
	stdErrors "errors"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/errgo.v1"
)

func TestErrCtx_RootCtx(t *testing.T) {
	t.Run("It should return a default context if the error is not wrapped with ErrCtx", func(t *testing.T) {
		// Given
		err := stdErrors.New("main error")
		err = errgo.Notef(err, "wrapping error in func2")

		// When
		rootCtx := RootCtx(err)

		// Then
		assert.Equal(t, context.Background(), rootCtx)
		assert.Contains(t, err.Error(), "wrapping error in func2")
	})

	t.Run("It should get the root context and contains fields from multiple wrapped error", func(t *testing.T) {
		// Given
		//nolint:all
		ctx := context.WithValue(context.Background(), "field0", "value0")
		err := funcThrowingError(ctx)
		err = Notef(ctx, err, "wrapping error in func2")
		err = Notef(ctx, err, "wrapping error in func3")
		err = Notef(ctx, err, "wrapping error in func4")

		// When
		rootCtx := RootCtx(err)
		assert.NotNil(t, rootCtx)

		// Then
		assert.Equal(t, "value0", rootCtx.Value("field0"))
		assert.Equal(t, "value1", rootCtx.Value("field1"))

		assert.Contains(t, err.Error(), "wrapping error in func2")
		assert.Contains(t, err.Error(), "wrapping error in func3")
		assert.Contains(t, err.Error(), "wrapping error in func4")
	})

	t.Run("It should get the root context and contains fields from function wrapping the error", func(t *testing.T) {
		// Given
		//nolint:all
		ctx := context.WithValue(context.Background(), "field0", "value0")
		err := funcWrappingAnError(ctx)
		err = Notef(ctx, err, "wrapping error in func3")
		err = Notef(ctx, err, "wrapping error in func4")

		// When
		rootCtx := RootCtx(err)
		assert.NotNil(t, rootCtx)

		// Then
		assert.Equal(t, "value0", rootCtx.Value("field0"))
		assert.Equal(t, "value1", rootCtx.Value("field1"))

		assert.Contains(t, err.Error(), "wrapping error from funcWrappingAnError")
		assert.Contains(t, err.Error(), "wrapping error in func3")
		assert.Contains(t, err.Error(), "wrapping error in func4")
	})

	t.Run("It should get the root context and contains first fields", func(t *testing.T) {
		// Given
		//nolint:all
		ctx := context.WithValue(context.Background(), "field0", "value0")
		// Simulate non ErrCtx error in middle of error path
		err := funcWrappingAnErrorWithoutErrCtx(ctx)
		err = Notef(ctx, err, "wrapping error in func2")
		err = Notef(ctx, err, "wrapping error in func3")

		// When
		rootCtx := RootCtx(err)
		assert.NotNil(t, rootCtx)

		// Then
		assert.Equal(t, "value0", rootCtx.Value("field0"))
		assert.Equal(t, "value1", rootCtx.Value("field1"))
		assert.Equal(t, "value2", rootCtx.Value("field2"))
		assert.Equal(t, "value3", rootCtx.Value("field3"))

		assert.Contains(t, err.Error(), "wrapping error from funcWrappingAnError")
		assert.Contains(t, err.Error(), "wrapping error from funcWrappingAnErrorWithoutErrCtx")
		assert.Contains(t, err.Error(), "wrapping error in func2")
		assert.Contains(t, err.Error(), "wrapping error in func3")
	})

	t.Run("It should get the root context and not contains first fields", func(t *testing.T) {
		// Given
		//nolint:all
		ctx := context.WithValue(context.Background(), "field0", "value0")
		err := funcThrowingError(ctx)
		assert.NotNil(t, err)
		// Simulate non returning error
		//nolint:all
		ctx = context.WithValue(ctx, "field2", "value2")
		err = Newf(ctx, "new error from func2")
		err = Notef(ctx, err, "wrapping error in func2")
		err = Notef(ctx, err, "wrapping error in func3")
		err = Notef(ctx, err, "wrapping error in func4")

		// When
		rootCtx := RootCtx(err)
		assert.NotNil(t, rootCtx)

		// Then
		assert.Equal(t, "value0", rootCtx.Value("field0"))
		assert.Equal(t, "value2", rootCtx.Value("field2"))
		assert.NotEqual(t, "value1", rootCtx.Value("field1"))

		assert.Contains(t, err.Error(), "new error from func2")
		assert.Contains(t, err.Error(), "wrapping error in func2")
		assert.Contains(t, err.Error(), "wrapping error in func3")
		assert.Contains(t, err.Error(), "wrapping error in func4")
	})
}

// funcThrowingError throw the main error
func funcThrowingError(ctx context.Context) error {
	//nolint:all
	ctx = context.WithValue(ctx, "field1", "value1")

	return Newf(ctx, "main error")
}

func funcWrappingAnError(ctx context.Context) error {
	//nolint:all
	ctx = context.WithValue(ctx, "field2", "value2")

	err := funcThrowingError(ctx)
	if err != nil {
		return Notef(ctx, err, "wrapping error from funcWrappingAnError")
	}
	return nil
}

func funcWrappingAnErrorWithoutErrCtx(ctx context.Context) error {
	//nolint:all
	ctx = context.WithValue(ctx, "field3", "value3")

	err := funcWrappingAnError(ctx)
	if err != nil {
		return errors.Wrap(err, "wrapping error from funcWrappingAnErrorWithoutErrCtx")
	}
	return nil
}
