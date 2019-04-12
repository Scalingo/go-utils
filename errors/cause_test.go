package errors

import (
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

	t.Run("given an error stack with errors.Wrapf", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{})
		err = errors.Wrapf(err, "biniou")

		assert.True(t, IsRootCause(err, &ValidationErrors{}))
		assert.False(t, IsRootCause(err, ValidationErrors{}))
	})
}

func Test_RootCause(t *testing.T) {
	t.Run("given an error stack with errgo.Mask", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": []string{"biniou"},
			},
		})
		err = errgo.Mask(err, errgo.Any)

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})

	t.Run("given an error stack with errgo.Notef", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": []string{"biniou"},
			},
		})
		err = errgo.Notef(err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})

	t.Run("given an error stack with errors.Wrap", func(t *testing.T) {
		var err error
		err = (&ValidationErrors{
			Errors: map[string][]string{
				"test": []string{"biniou"},
			},
		})
		err = errors.Wrap(err, "pouet")

		assert.Equal(t, "test=biniou", RootCause(err).Error())
	})
}
