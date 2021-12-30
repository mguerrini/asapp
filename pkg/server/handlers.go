package server
import (
	"context"
	"github.com/challenge/pkg/domain/security"
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/server"
	"github.com/challenge/pkg/services"
	"net/http"
)

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
	next.Handle(ctx, w, r)
}

func (ah *RequestHandler) ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.HttpInterceptor)  {
	next.Handle(ctx, w, r)
}



