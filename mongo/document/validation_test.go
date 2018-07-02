package document

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

type DummyDocument struct {
	Base
	FieldAErrors int
	FieldBErrors int
}

func (d *DummyDocument) Validate(ctx context.Context) *ValidationErrors {
	err := NewValidationErrorsBuilder()

	for i := 0; i < d.FieldAErrors; i++ {
		err.Set("a", "test")
	}

	for i := 0; i < d.FieldBErrors; i++ {
		err.Set("b", "test")
	}

	return err.Build()

}

func TestValidation(t *testing.T) {
	examples := map[string]struct {
		ExpectedError error
		Document      *DummyDocument
	}{
		"no errors": {
			Document:      &DummyDocument{},
			ExpectedError: nil,
		},
		"with some errors": {
			Document: &DummyDocument{
				FieldAErrors: 1,
				FieldBErrors: 2,
			},
			ExpectedError: &ValidationErrors{
				Errors: map[string][]string{
					"a": []string{"test"},
					"b": []string{"test", "test"},
				},
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		for name, example := range examples {
			t.Run(name, func(t *testing.T) {
				d := example.Document
				err := Create(context.Background(), "test", d)

				assert.Equal(t, example.ExpectedError, err)
			})
		}
	})

	t.Run("save", func(t *testing.T) {
		for name, example := range examples {
			t.Run(name, func(t *testing.T) {
				d := example.Document
				err := Save(context.Background(), "test", d)

				assert.Equal(t, example.ExpectedError, err)
			})
		}
	})

	t.Run("update", func(t *testing.T) {
		for name, example := range examples {
			t.Run(name, func(t *testing.T) {
				d := example.Document
				err := Update(context.Background(), "test", bson.M{}, d)

				assert.Equal(t, example.ExpectedError, err)
			})
		}
	})
}

func TestValidationErrors_Error(t *testing.T) {
	t.Run("it should return a specific error message if nil", func(t *testing.T) {
		var err *ValidationErrors
		require.Equal(t, "empty validation error", err.Error())
	})
}
