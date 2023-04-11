package errors

import (
	"errors"
)

// Is checks if any error of the stack matches the error value expectedError
// API machting the standard library but allowing to wrap errors with ErrCtx + errgo or pkg/errors
func Is(receivedErr, expectedError error) bool {
	if errors.Is(receivedErr, expectedError) {
		return true
	}
	for receivedErr != nil {
		receivedErr = Unwrap(receivedErr)
		if errors.Is(receivedErr, expectedError) {
			return true
		}
	}
	return false
}

// As checks if any error of the stack matches the expectedType
// API machting the standard library but allowing to wrap errors with ErrCtx + errgo or pkg/errors
func As(receivedErr error, expectedType any) bool {
	if errors.As(receivedErr, expectedType) {
		return true
	}
	for receivedErr != nil {
		receivedErr = Unwrap(receivedErr)
		if errors.As(receivedErr, expectedType) {
			return true
		}
	}
	return false
}

// Unwrap tries to unwrap `err`, getting the wrapped error or nil.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
