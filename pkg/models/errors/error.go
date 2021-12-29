package errors

import (
	"errors"
	"net/http"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int   `json:"status"`
	Err  error `json:"error"`
}

func NewBadRequest(err error) *StatusError {
	return NewStatusError(http.StatusBadRequest, err)
}

func NewBadRequestMsg(msg string) *StatusError {
	return NewStatusError(http.StatusNotFound, errors.New(msg))
}

func NewNotFound(err error) *StatusError {
	return NewStatusError(http.StatusBadRequest, err)
}

func NewNotFoundMsg(msg string) *StatusError {
	return NewStatusError(http.StatusNotFound, errors.New(msg))
}

func NewInternalServerError(err error) *StatusError {
	return NewStatusError(http.StatusInternalServerError, err)
}

func NewInternalServerErrorMsg(msg string ) *StatusError {
	return NewStatusError(http.StatusInternalServerError, errors.New(msg))
}


func NewStatusError(status int, err error ) *StatusError {
	return &StatusError{
		Err:  err,
		Code: status,
	}
}

func NewStatusErrorMsg(status int, msg string) *StatusError {
	return &StatusError{
		Err:  errors.New(msg),
		Code: status,
	}
}

func (e *StatusError) Status() int {
	return e.Code
}

func (e *StatusError) Error() string {
	return e.Err.Error()
}
