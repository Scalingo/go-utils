package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidationErrorsBuilder_Merge(t *testing.T) {
	cases := map[string]struct {
		Builder  *ValidationErrorsBuilder
		Error    *ValidationErrors
		Expected *ValidationErrorsBuilder
	}{
		"merging nil is a no-op": {
			Builder:  NewValidationErrorsBuilder().Set("field", "invalid"),
			Error:    NewValidationErrorsBuilder().Build(),
			Expected: NewValidationErrorsBuilder().Set("field", "invalid"),
		},
		"merging should add them to the builder": {
			Builder:  NewValidationErrorsBuilder().Set("field", "invalid"),
			Error:    NewValidationErrorsBuilder().Set("f1", "err").Set("f2", "err").Build(),
			Expected: NewValidationErrorsBuilder().Set("field", "invalid").Set("f1", "err").Set("f2", "err"),
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			require.Equal(t, c.Expected, c.Builder.Merge(c.Error))
		})
	}
}

func TestValidationErrorsBuilder_MergeWithPrefix(t *testing.T) {
	cases := map[string]struct {
		Builder  *ValidationErrorsBuilder
		Error    *ValidationErrors
		Expected *ValidationErrorsBuilder
		Prefix   string
	}{
		"merging nil is a no-op": {
			Builder:  NewValidationErrorsBuilder().Set("field", "invalid"),
			Error:    NewValidationErrorsBuilder().Build(),
			Expected: NewValidationErrorsBuilder().Set("field", "invalid"),
		},
		"merging should add them to the builder with the prefix + '_'": {
			Prefix:   "a",
			Builder:  NewValidationErrorsBuilder().Set("field", "invalid"),
			Error:    NewValidationErrorsBuilder().Set("f1", "err").Set("f2", "err").Build(),
			Expected: NewValidationErrorsBuilder().Set("field", "invalid").Set("a_f1", "err").Set("a_f2", "err"),
		},
		"merging a fields should add them to the builder with the prefix without adding '_' if present": {
			Prefix:   "a_",
			Builder:  NewValidationErrorsBuilder().Set("field", "invalid"),
			Error:    NewValidationErrorsBuilder().Set("f1", "err").Set("f2", "err").Build(),
			Expected: NewValidationErrorsBuilder().Set("field", "invalid").Set("a_f1", "err").Set("a_f2", "err"),
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			require.Equal(t, c.Expected, c.Builder.MergeWithPrefix(c.Prefix, c.Error))
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	cases := map[string]struct {
		Errors         ValidationErrors
		expectedErrors []string
	}{
		"should return a string with one error in it": {
			Errors: ValidationErrors{
				Errors: map[string][]string{
					"name": {"invalid name"},
				},
			},
			expectedErrors: []string{"name=invalid name"},
		},
		"should return a string with multiple errors in it with the same field name": {
			Errors: ValidationErrors{
				Errors: map[string][]string{
					"name": {"invalid name", "should contains alphanumeric characters"},
				},
			},
			expectedErrors: []string{"name=invalid name", "should contains alphanumeric characters"},
		},
		"should return a string with multiple errors in it with multiple field name": {
			Errors: ValidationErrors{
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
				require.Contains(t, c.Errors.Error(), expectedError)
			}
		})
	}
}
