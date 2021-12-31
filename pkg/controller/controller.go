package controller

import (
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/storage"
	"github.com/challenge/pkg/repository"
	"github.com/challenge/pkg/services"
)

// Handler provides the interface to handle different requests
type Handler struct {
	healthRepository repository.IHealthRepository

	msgService   *services.MessageServices
	userServices *services.UserServices
	authServices *services.AuthServices
}

func NewController() *Handler {
	sessionName, err := config.ConfigurationSingleton().GetString("root.startup.session_name")

	if err != nil {
		panic("There no session name configured (path: root.startup.session_name). " + err.Error())
	}

	 storage.DBManagerSingleton().CreateConnection(sessionName)

	return &Handler{
		healthRepository: repository.RepositoryFactory().CreateHealthRepository(sessionName),

		msgService:   services.NewMessageServices(sessionName),
		userServices: services.NewUserServices(sessionName),
		authServices: services.NewAuthServices(sessionName),
	}
}