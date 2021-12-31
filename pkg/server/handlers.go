package server
import (
	"context"
	"github.com/challenge/pkg/domain/security"
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/server"
	"github.com/challenge/pkg/modules/storage"
	"github.com/challenge/pkg/services"
	"net/http"
)

type statusCaptureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusCaptureResponseWriter(w http.ResponseWriter) *statusCaptureResponseWriter {
	return &statusCaptureResponseWriter{w, http.StatusOK}
}

func (lrw *statusCaptureResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}


type RequestHandler struct {
	authService *services.AuthServices
}

func NewRequestHandler() *RequestHandler {
	sessionName, err := config.ConfigurationSingleton().GetString("root.startup.session_name")

	if err != nil {
		panic("There no session name configured (path: root.startup.session_name). " + err.Error())
	}

	authServ := services.NewAuthServices(sessionName)

	return &RequestHandler{
		authService: authServ,
	}
}

func (ah *RequestHandler) ValidateTokenHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.HttpInterceptor)  {
	//get token
	token := r.Header.Get("Authorization")
	status := ah.authService.ValidateToken(ctx, token)

	switch status {
	case security.SecurityTokenStatus_OK:
		next.Handle(ctx, w, r)
		return
	case security.SecurityTokenStatus_Expired:
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return

	case security.SecurityTokenStatus_Invalid:
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
}



func (ah *RequestHandler) TransactionScopeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.HttpInterceptor)  {
	newCtx, scope := storage.TransactionalScopeManagerSingleton().Begin(ctx, nil)

	newW := NewStatusCaptureResponseWriter(w)

	defer func() {
		if newW.statusCode == http.StatusOK {
			scope.Commit(newCtx)
		} else {
			scope.Rollback(newCtx)
		}
	}()

	next.Handle(newCtx, newW, r)
}

func (ah *RequestHandler) ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.HttpInterceptor)  {
	next.Handle(ctx, w, r)
}



