package errors

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	errgo "gopkg.in/errgo.v1"
)

func Test_IsRootCause(t *testing.T) {
	t.Run("given an error stack with errgo.Notef", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errgo.Notef(err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})

	t.Run("given an error stack with errgo.NoteMask", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errgo.NoteMask(err, "biniou")

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

	t.Run("given an error stack with NoteMask from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = NoteMask(context.Background(), err, "biniou")

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

	t.Run("given an error stack with errgo.NoteMask", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = errgo.NoteMask(err, "pouet")

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

	t.Run("given an error stack with NoteMask from ErrCtx", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": {"biniou"},
			},
		})
		err = NoteMask(context.Background(), err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})
}
