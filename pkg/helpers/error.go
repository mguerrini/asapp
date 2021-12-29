package helpers

import (
	apierrors "github.com/challenge/pkg/models/errors"
	"net/http"
)

func GetStatusCode(err error) int {
	switch e := err.(type) {
		case apierrors.Error:
			return e.Status()
	 	default:
	 		return http.StatusInternalServerError
	}
}

func GetStatusCodeOr(err error, status int) int {
	switch e := err.(type) {
	case apierrors.Error:
		return e.Status()
	default:
		return status
	}
}
