package errors

import "gopkg.in/errgo.v1"

func ErrgoRoot(err error) error {
	for {
		e, ok := err.(*errgo.Err)
		if !ok {
			return err
		}
		if e.Underlying() == nil {
			return err
		}
		err = e.Underlying()
	}
}

type UserFacingError struct {
	UserMessage    string
	TechnicalError error
}

func (err *UserFacingError) Error() string {
	return err.UserMessage
}

func WrapUseMessageAroundError(err error, wrappingMessage string) error {
	userFacingError := &UserFacingError{
		UserMessage:    wrappingMessage,
		TechnicalError: err,
	}
	return userFacingError
}
