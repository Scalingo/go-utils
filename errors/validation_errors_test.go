package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationErrorsBuilder_Build(t *testing.T) {
	t.Run("empty validation error is a nil error", func(t *testing.T) {
		err := NewValidationErrorsBuilder().Build()
		assert.NoError(t, err)
	})

	t.Run("Build must return a correct validation error", func(t *testing.T) {
		err := NewValidationErrorsBuilder().Set("field", "invalid").Build()

		var verr *ValidationErrors
		ok := errors.As(err, &verr)
		require.True(t, ok)

		assert.EqualError(t, err, "field=invalid")
	})
}

func TestValidationErrorsBuilder_Merge(t *testing.T) {
	cases := map[string]struct {
		validationErrorsBuilder *ValidationErrorsBuilder
		validationErrorsToMerge error
		expectedBuilder         *ValidationErrorsBuilder
	}{
		"merging nil is a no-op": {
			validationErrorsBuilder: NewValidationErrorsBuilder().Set("field", "invalid"),
			validationErrorsToMerge: NewValidationErrorsBuilder().Build(),
			expectedBuilder:         NewValidationErrorsBuilder().Set("field", "invalid"),
		},
		"merging should add them to the builder": {
			validationErrorsBuilder: NewValidationErrorsBuilder().Set("field", "invalid"),
			validationErrorsToMerge: NewValidationErrorsBuilder().Set("f1", "err").Set("f2", "err").Build(),
			expectedBuilder:         NewValidationErrorsBuilder().Set("field", "invalid").Set("f1", "err").Set("f2", "err"),
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			var verr *ValidationErrors
			if c.validationErrorsToMerge != nil {
				ok := errors.As(c.validationErrorsToMerge, &verr)
				require.True(t, ok)
			}

			mergedError := c.validationErrorsBuilder.Merge(verr)
			assert.Equal(t, c.expectedBuilder, mergedError)
		})
	}
}

func TestValidationErrorsBuilder_MergeWithPrefix(t *testing.T) {
	cases := map[string]struct {
		prefix                  string
		validationErrorsBuilder *ValidationErrorsBuilder
		validationErrorsToMerge error
		expectedBuilder         *ValidationErrorsBuilder
	}{
		"merging nil is a no-op": {
			validationErrorsBuilder: NewValidationErrorsBuilder().Set("field", "invalid"),
			validationErrorsToMerge: NewValidationErrorsBuilder().Build(),
			expectedBuilder:         NewValidationErrorsBuilder().Set("field", "invalid"),
		},
		"merging should add them to the builder with the prefix + '_'": {
			prefix:                  "a",
			validationErrorsBuilder: NewValidationErrorsBuilder().Set("field", "invalid"),
			validationErrorsToMerge: NewValidationErrorsBuilder().Set("f1", "err").Set("f2", "err").Build(),
			expectedBuilder:         NewValidationErrorsBuilder().Set("field", "invalid").Set("a_f1", "err").Set("a_f2", "err"),
		},
		"merging a fields should add them to the builder with the prefix without adding '_' if present": {
			prefix:                  "a_",
			validationErrorsBuilder: NewValidationErrorsBuilder().Set("field", "invalid"),
			validationErrorsToMerge: NewValidationErrorsBuilder().Set("f1", "err").Set("f2", "err").Build(),
			expectedBuilder:         NewValidationErrorsBuilder().Set("field", "invalid").Set("a_f1", "err").Set("a_f2", "err"),
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			var verr *ValidationErrors
			if c.validationErrorsToMerge != nil {
				ok := errors.As(c.validationErrorsToMerge, &verr)
				require.True(t, ok)
			}

			mergedError := c.validationErrorsBuilder.MergeWithPrefix(c.prefix, verr)
			require.Equal(t, c.expectedBuilder, mergedError)
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	cases := map[string]struct {
		validationErrors ValidationErrors
		expectedErrors   []string
	}{
		"should return a string with one error in it": {
			validationErrors: ValidationErrors{
				Errors: map[string][]string{
					"name": {"invalid name"},
				},
			},
			expectedErrors: []string{"name=invalid name"},
		},
		"should return a string with multiple errors in it with the same field name": {
			validationErrors: ValidationErrors{
				Errors: map[string][]string{
					"name": {"invalid name", "should contains alphanumeric characters"},
				},
			},
			expectedErrors: []string{"name=invalid name", "should contains alphanumeric characters"},
		},
		"should return a string with multiple errors in it with multiple field name": {
			validationErrors: ValidationErrors{
				Errors: map[string][]string{
					"name": {"invalid name", "should contains alphanumeric characters"},
					"type": {"invalid type", "type not exist"},
				},
			},
			expectedErrors: []string{"name=invalid name, should contains alphanumeric characters", "type=invalid type, type not exist"},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			for _, expectedError := range c.expectedErrors {
				require.Contains(t, c.validationErrors.Error(), expectedError)
			}
		})
	}
}
