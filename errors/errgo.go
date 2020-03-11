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

type UserFacingError interface {
	error
	UserFacingError() string
	TechnicalError() error
}

type GenericUserFacingError struct {
	technicalError  error
	userFacingError string
}

func (err GenericUserFacingError) Error() string {
	return err.technicalError.Error()
}

func (err GenericUserFacingError) TechnicalError() error {
	return err.technicalError
}

func (err GenericUserFacingError) UserFacingError() string {
	return err.userFacingError
}

func WrapUseMessageAroundError(err error, msg string) UserFacingError {
	return GenericUserFacingError{
		technicalError:  err,
		userFacingError: msg,
	}
}

func IsUserFacingError(err error) bool {
	if err == nil {
		return false
	}
	_, is := err.(UserFacingError)
	return is
}
