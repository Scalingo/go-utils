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

	t.Run("given an error stack with Notef from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = Notef(context.Background(), err, "biniou")

		assert.True(t, As(err, &expectedErrorType))
		assert.False(t, As(err, &unexpectedErrorType))
	})

	t.Run("given an error in the middle of the stack stack with Notef from ErrCtx", func(t *testing.T) {
		var err error
		err = io.EOF
		err = &customError{WrappedError: err, CustomValue: "value"}
		err = Notef(context.Background(), err, "biniou")

		var expectedErrorType *customError
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
		err := Notef(context.Background(), expectedError, "pouet")
		err = errgo.Notef(err, "pouet")
		err = errors.Wrap(err, "pouet")

		assert.True(t, Is(err, expectedError))
	})
}

func Test_IsRootCause(t *testing.T) {
	t.Run("given an error stack with errgo.Notef", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errgo.Notef(err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})

	t.Run("given an error stack with errors.Wrap", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errors.Wrap(err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})

	t.Run("given an error stack with errors.Wrapf", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errors.Wrapf(err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})

	t.Run("given an error stack with Wrap from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = Wrap(context.Background(), err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})

	t.Run("given an error stack with Wrapf from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = Wrapf(context.Background(), err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})

	t.Run("given an error stack with Notef from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = Notef(context.Background(), err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})

}

func Test_RootCause(t *testing.T) {
	t.Run("given an error stack with errgo.Mask", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = errgo.Mask(err, errgo.Any)

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})

	t.Run("given an error stack with errgo.Notef", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = errgo.Notef(err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})

	t.Run("given an error stack with errors.Wrap", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = errors.Wrap(err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})

	t.Run("given an error stack with Wrap from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = Wrap(context.Background(), err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})

	t.Run("given an error stack with Wrapf from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = Wrapf(context.Background(), err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})

	t.Run("given an error stack with Notef from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = Notef(context.Background(), err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
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

	t.Run("given an error stack with Notef from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = Notef(context.Background(), err, "pouet")

		assert.Equal(t, "pouet: test=biniou", UnwrapError(err).Error())
	})
	t.Run("given an error nil", func(t *testing.T) {
		var err error
		assert.Nil(t, UnwrapError(err))
	})
}
