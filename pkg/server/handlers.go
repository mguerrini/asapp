package server
import (
	"context"
	"net/http"
)

type RequestHandler struct {

}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func (ah *RequestHandler) ValidateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	if false {
		// TODO: validate token
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return false
	}

	return true
}

func (ah *RequestHandler) TransactionScopeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	if false {
		// TODO: validate token
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return false
	}

	return true
}

func (ah *RequestHandler) ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	return true
}



