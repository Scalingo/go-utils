package errors

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	errgo "gopkg.in/errgo.v1"
)

type customError struct {
	WrappedError error
	CustomValue  string
}

func (err *customError) Error() string {
	return "custom error " + err.CustomValue
}

func (err *customError) Unwrap() error {
	return err.WrappedError
}

func Test_As(t *testing.T) {
	var expectedErrorType *ValidationErrors
	var unexpectedErrorType *net.OpError
	t.Run("given an error stack with errgo.Notef", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errgo.Notef(err, "biniou")

		assert.True(t, As(err, &expectedErrorType))
		assert.False(t, As(err, &unexpectedErrorType))
	})

	t.Run("given an error stack with errors.Wrap", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errors.Wrap(err, "biniou")

		assert.True(t, As(err, &expectedErrorType))
		assert.False(t, As(err, &unexpectedErrorType))
	})

	t.Run("given an error stack with errors.Wrapf", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errors.Wrapf(err, "biniou")

		assert.True(t, As(err, &expectedErrorType))
		assert.False(t, As(err, &unexpectedErrorType))
	})

	t.Run("given an error stack with Wrap from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = Wrap(context.Background(), err, "biniou")

		assert.True(t, As(err, &expectedErrorType))
		assert.False(t, As(err, &unexpectedErrorType))
	})

	t.Run("given an error stack with Wrapf from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = Wrapf(context.Background(), err, "biniou")

		assert.True(t, As(err, &expectedErrorType))
		assert.False(t, As(err, &unexpectedErrorType))
	})
}

func Test_Is(t *testing.T) {
	t.Run("given an error stack with errgo.Mask", func(t *testing.T) {
		expectedError := io.EOF
		err := errgo.Mask(expectedError, errgo.Any)

		assert.True(t, Is(err, expectedError))
	})

	t.Run("given an error stack with errgo.Notef", func(t *testing.T) {
		expectedError := io.EOF
		err := errgo.Notef(expectedError, "pouet")

		assert.True(t, Is(err, expectedError))
	})

	t.Run("given an error stack with errors.Wrap", func(t *testing.T) {
		expectedError := io.EOF
		err := errors.Wrap(expectedError, "pouet")

		assert.True(t, Is(err, expectedError))
	})

	t.Run("given an error stack with Wrap from ErrCtx", func(t *testing.T) {
		expectedError := io.EOF
		err := Wrap(context.Background(), expectedError, "pouet")

		assert.True(t, Is(err, expectedError))
	})

	t.Run("given an error stack with Wrapf from ErrCtx", func(t *testing.T) {
		expectedError := io.EOF
		err := Wrapf(context.Background(), expectedError, "pouet")

		assert.True(t, Is(err, expectedError))
	})

	t.Run("given an error stack with mixed types", func(t *testing.T) {
		expectedError := io.EOF

		err := Wrapf(t.Context(), io.EOF, "pouet")
		err = errgo.Notef(err, "pouet")
		err = errors.Wrap(err, "pouet")

		assert.True(t, Is(err, expectedError))
	})
}

func Test_UnwrapError(t *testing.T) {
	t.Run("given an error stack with errgo.Mask", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = errgo.Mask(err, errgo.Any)

		assert.Equal(t, "test=biniou", UnwrapError(err).Error())
	})

	t.Run("given an error stack multiple times with errors.Wrap", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = errors.Wrap(err, "pouet")
		err = errors.Wrap(err, "pouet")
		err = errors.Wrap(err, "pouet")
		err = errors.Wrap(err, "pouet")
		err = errors.Wrap(err, "pouet")

		var lastErr error
		for unwrappedErr := err; unwrappedErr != nil; unwrappedErr = UnwrapError(unwrappedErr) {
			lastErr = unwrappedErr
		}

		assert.Equal(t, "test=biniou", lastErr.Error())
	})

	t.Run("given a nil error", func(t *testing.T) {
		var err error
		assert.Nil(t, UnwrapError(err))
	})
}
