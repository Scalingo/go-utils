package document

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

type DummyDocument struct {
	Base
	FieldAErrors  int
	FieldBErrors  int
	InternalError bool
}

func (d *DummyDocument) Validate(ctx context.Context) (*ValidationErrors, error) {
	err := NewValidationErrorsBuilder()

	if d.InternalError {
		return nil, errors.New("Internal error")
	}

	for i := 0; i < d.FieldAErrors; i++ {
		err.Set("a", "test")
	}

	for i := 0; i < d.FieldBErrors; i++ {
		err.Set("b", "test")
	}

	return err.Build(), nil
}

func TestValidation(t *testing.T) {
	examples := map[string]struct {
		ExpectedError           error
		ExpectedValidationError error
		Document                *DummyDocument
	}{
		"no errors": {
			Document:                &DummyDocument{},
			ExpectedError:           nil,
			ExpectedValidationError: nil,
		},
		"with some validation errors": {
			Document: &DummyDocument{
				FieldAErrors: 1,
				FieldBErrors: 2,
			},
			ExpectedError: nil,
			ExpectedValidationError: &ValidationErrors{
				Errors: map[string][]string{
					"a": []string{"test"},
					"b": []string{"test", "test"},
				},
			},
		},
		"with internal error": {
			Document: &DummyDocument{
				InternalError: true,
			},
			ExpectedValidationError: nil,
			ExpectedError:           errors.New("Internal error"),
		},
	}

	t.Run("create", func(t *testing.T) {
		for name, example := range examples {
			t.Run(name, func(t *testing.T) {
				d := example.Document
				err := Create(context.Background(), "test", d)

				if example.ExpectedError == nil && example.ExpectedValidationError == nil {
					assert.NoError(t, err)
				}
				if example.ExpectedError != nil {
					assert.Equal(t, example.ExpectedError, err)
				}
				if example.ExpectedValidationError != nil {
					assert.Equal(t, example.ExpectedValidationError, err)
				}
			})
		}
	})

	t.Run("save", func(t *testing.T) {
		for name, example := range examples {
			t.Run(name, func(t *testing.T) {
				d := example.Document
				err := Save(context.Background(), "test", d)

				if example.ExpectedError == nil && example.ExpectedValidationError == nil {
					assert.NoError(t, err)
				}
				if example.ExpectedError != nil {
					assert.Equal(t, example.ExpectedError, err)
				}
				if example.ExpectedValidationError != nil {
					assert.Equal(t, example.ExpectedValidationError, err)
				}
			})
		}
	})

	t.Run("update", func(t *testing.T) {
		for name, example := range examples {
			t.Run(name, func(t *testing.T) {
				d := example.Document
				err := Update(context.Background(), "test", bson.M{}, d)

				if example.ExpectedError == nil && example.ExpectedValidationError == nil {
					assert.NoError(t, err)
				}
				if example.ExpectedError != nil {
					assert.Equal(t, example.ExpectedError, err)
				}
				if example.ExpectedValidationError != nil {
					assert.Equal(t, example.ExpectedValidationError, err)
				}
			})
		}
	})
}
