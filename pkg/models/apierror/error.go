package apierrors

import (
	"errors"
	"fmt"
	"net/http"
)

type GenericError interface {
	GetStatus() int
	Error() string
	Cause() error
	ErrorClass() string
	AsString() string
}

type ApiError struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Err     string `json:"error"`
	Message string `json:"message"`
}


func NewBadRequest(err error, message string) *ApiError {
	return NewApiError(err, message, http.StatusBadRequest)
}


func NewInternalServerError(err error, message string) *ApiError {
	return NewApiError(err, message, http.StatusInternalServerError)
}

func NewApiError(err error, message string, status int) *ApiError {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	} else {
		errMsg = message
	}
	return &ApiError{
		Message: message,
		Err:     errMsg,
		Status:  http.StatusText(status),
		Code:    status,
	}
}


func (e *ApiError) GetMessage() string {
	return e.Message
}

func (e *ApiError) GetStatusCode() int {
	return e.Code
}

func (e *ApiError) GetError() string {
	return e.Err
}

func (e *ApiError) AsString() string {
	if len(e.Err) == 0 {
		if len(e.GetMessage()) == 0 {
			return fmt.Sprintf("[%d - %s]", e.Code, e.Status)
		} else {
			return fmt.Sprintf("[%d - %s] message: %s", e.Code, e.Status, e.Message)
		}
	} else {
		if len(e.GetMessage()) == 0 {
			return fmt.Sprintf("[%d - %s] error: %s", e.Code, e.Status, e.Err)
		} else {
			return fmt.Sprintf("[%d - %s] message: %s, error: %s", e.Code, e.Status, e.Message, e.Err)
		}
	}
}

/*************************************/
/*      GenericError interface       */
/*************************************/
func (e *ApiError) GetStatus() int {
	return e.Code
}

func (e *ApiError) Error() string {
	return e.Err
}

func (e *ApiError) Cause() error {
	if len(e.Err) == 0 {
		return errors.New(e.Message)
	} else {
		return errors.New(e.Err)
	}
}

func (e *ApiError) ErrorClass() string {
	return e.Status
}

/*************************************/
