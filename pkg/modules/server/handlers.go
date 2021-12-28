package server

import (
	"context"
	"github.com/challenge/pkg/helpers"
	"net/http"
)

type HttpHandler struct {
	Root     *StageHttpHandler
	GetRoot  *StageHttpHandler
	PostRoot *StageHttpHandler
}


//Return a boolean that indicates if can continue with the execution
type HttpCancelHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool

type HttpHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)


type StageHttpHandler struct {
	Handler HttpCancelHandlerFunc
	Next    *StageHttpHandler
}


func NewhttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) AddCancelHandlerFunc(handle HttpCancelHandlerFunc) *HttpHandler {
	handler := &StageHttpHandler{
		Handler: handle,
		Next:    nil,
	}

	return h.AddHandler(handler)
}

func (h *HttpHandler) AddHandlerFunc(handle HttpHandlerFunc) *HttpHandler {

	adapter := func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
		handle(ctx, w, r)
		return true
	}

	handler := &StageHttpHandler{
		Handler: adapter,
		Next:    nil,
	}

	return h.AddHandler(handler)
}

func (h *HttpHandler) AddHandler( handle *StageHttpHandler) *HttpHandler {
	if h.Root == nil {
		h.Root = handle
		return h
	}

	h.Root.AddNext(handle)
	return h
}


func (h *HttpHandler) AddCancelMethodHandlerFunc(httpMethod string, handle HttpCancelHandlerFunc) *HttpHandler {
	handler := &StageHttpHandler{
		Handler: handle,
		Next:    nil,
	}

	return h.AddMethodHandler(httpMethod, handler)
}

func (h *HttpHandler) AddMethodHandlerFunc(httpMethod string, handle HttpHandlerFunc) *HttpHandler {

	adapter := func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
		handle(ctx, w, r)
		return true
	}

	handler := &StageHttpHandler{
		Handler: adapter,
		Next:    nil,
	}

	return h.AddMethodHandler(httpMethod, handler)
}

func (h *HttpHandler) AddMethodHandler(httpMethod string, handle *StageHttpHandler) *HttpHandler {
	if httpMethod == http.MethodGet {
		if h.GetRoot == nil {
			h.GetRoot = handle
			return h
		}

		h.GetRoot.AddNext(handle)
	} else if httpMethod == http.MethodPost {
		if h.PostRoot == nil {
			h.PostRoot = handle
			return h
		}

		h.PostRoot.AddNext(handle)
	} else {
		h.AddHandler(handle)
	}
	return h
}



func (h *HttpHandler) Handle (w http.ResponseWriter, r *http.Request) {
	//default error handler
	defer helpers.Recover()

	method := r.Method

	if h.GetRoot != nil || h.PostRoot != nil {
		if h.GetRoot != nil && method == http.MethodGet {
			h.GetRoot.Handle(r.Context(), w, r)
		} else if h.PostRoot != nil && method == http.MethodPost {
			h.PostRoot.Handle(r.Context(), w, r)
		} else {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return
		}
	} else {
		if h.Root == nil {
			return
		}

		h.Root.Handle(r.Context(), w, r)
	}
}

func (h *StageHttpHandler) Handle (ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cont := h.Handler(ctx, w, r)

	if cont && h.Next != nil {
		h.Next.Handle(ctx, w, r)
	}
}

func (h *StageHttpHandler) AddNext (handle *StageHttpHandler) {
	if h.Next == nil {
		h.Next = handle
		return
	}

	h.Next.AddNext(handle)
}

