package server

import (
	"context"
	"fmt"
	"github.com/challenge/pkg/modules/logger"
	"net/http"
	"runtime/debug"
)

type HttpHandler struct {
	GetRoot  *stageHttpHandler
	PostRoot *stageHttpHandler
}

type HttpInterceptor interface {
	Handle(ctx context.Context, w http.ResponseWriter, r *http.Request)
}


type HttpHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

type HttpInterceptorHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, next HttpInterceptor)


type stageHttpHandler struct {
	Handler HttpInterceptorHandlerFunc
	Next    *stageHttpHandler
}


func NewhttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) AddInterceptorFunc(httpMethod string, handle HttpInterceptorHandlerFunc) *HttpHandler {
	handler := &stageHttpHandler{
		Handler: handle,
		Next:    nil,
	}

	return h.AddHandler(httpMethod, handler)
}

func (h *HttpHandler) AddHandlerFunc(httpMethod string, handle HttpHandlerFunc) *HttpHandler {

	adapter := func(ctx context.Context, w http.ResponseWriter, r *http.Request, next HttpInterceptor)  {
		handle(ctx, w, r)
	}

	handler := &stageHttpHandler{
		Handler: adapter,
		Next:    nil,
	}

	return h.AddHandler(httpMethod, handler)
}

func (h *HttpHandler) AddHandler(httpMethod string, handle *stageHttpHandler) *HttpHandler {
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
		panic("Invalid method")
	}
	return h
}


func (h *HttpHandler) Handle (w http.ResponseWriter, r *http.Request) {
	err := h.handle(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HttpHandler) handle (w http.ResponseWriter, r *http.Request) (err error) {
	//default error handler
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[Custom Recovery] panic recovered: %s %s", fmt.Errorf("%s", r), debug.Stack(), r)
			err = r.(error)
		}
	}()

	method := r.Method

	if h.GetRoot != nil || h.PostRoot != nil {
		if h.GetRoot != nil && method == http.MethodGet {
			h.GetRoot.doHandle(r.Context(), w, r)
		} else if h.PostRoot != nil && method == http.MethodPost {
			h.PostRoot.doHandle(r.Context(), w, r)
		} else {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return nil
		}
	} else {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return nil
	}

	return nil
}


func (h *stageHttpHandler) doHandle (ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.Handler(ctx, w, r, h)
}

func (h *stageHttpHandler) Handle (ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if h.Next == nil {
		return
	}

	h.Next.doHandle(ctx, w, r)
}

func (h *stageHttpHandler) AddNext (handle *stageHttpHandler) {
	if h.Next == nil {
		h.Next = handle
		return
	}

	h.Next.AddNext(handle)
}

