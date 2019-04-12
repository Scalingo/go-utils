package document

import (
	"github.com/Scalingo/go-utils/errors"
)

// ValidationErrors is a type alias of errors.ValidationErrors. It is defined to
// keep retro-compatibility
type ValidationErrors = errors.ValidationErrors

// ValidationErrorsBuilder is a type alias of errors.ValidationErrorsBuilder.
// It is defined to keep retro-compatibility
type ValidationErrorsBuilder = errors.ValidationErrorsBuilder

// NewValidationErrors return an empty ValidationErrors struct
func NewValidationErrorsBuilder() *ValidationErrorsBuilder {
	return errors.NewValidationErrorsBuilder()
}
