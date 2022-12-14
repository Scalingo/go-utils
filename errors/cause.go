package errors

import (
	"reflect"
)

// IsRootCause return true if the cause of the given error is the same type as
// mytype.
// This function takes the cause of an error if the errors stack has been
// wrapped with errors.Wrapf or errgo.Notef or errgo.NoteMask or errgo.Mask.
//
// Example:
//
//	errors.IsRootCause(err, &ValidationErrors{})
func IsRootCause(err error, mytype interface{}) bool {
	t := reflect.TypeOf(mytype)
	errCause := errorCause(err)
	errRoot := errgoRoot(err)
	return reflect.TypeOf(errCause) == t || reflect.TypeOf(errRoot) == t
}

// RootCause returns the cause of an errors stack, whatever the method they used
// to be stacked: either errgo.Notef or errors.Wrapf.
func RootCause(err error) error {
	errCause := errorCause(err)
	if errCause == nil {
		errCause = errgoRoot(err)
	}
	return errCause
}
