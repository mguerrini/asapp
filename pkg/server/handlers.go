package server
import (
	"context"
	"github.com/challenge/pkg/modules/server"
	"net/http"
)

type RequestHandler struct {

}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func (ah *RequestHandler) ValidateUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.HttpInterceptor)  {
	if false {
		// TODO: validate token
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	next.Handle(ctx, w, r)
}

func (ah *RequestHandler) TransactionScopeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.HttpInterceptor)  {
	next.Handle(ctx, w, r)
}

func (ah *RequestHandler) ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.HttpInterceptor)  {
	next.Handle(ctx, w, r)
}



