package server

import (
	"github.com/challenge/pkg/controller"
	"github.com/challenge/pkg/modules/server"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort = "8080"
	CheckEndpoint = "/check"
	UsersEndpoint = "/users"
	LoginEndpoint = "/login"
	MessagesEndpoint = "/messages"
)

func StartHttpServer(wg *sync.WaitGroup) *http.Server {
	StartServer() //initialize components

	//create server
	srv := &http.Server{Addr: ":" + ServerPort}

	//create handler
	controller := controller.NewController() // controller.Handler{}
	handler := NewRequestHandler()

	// Configure endpoints
	// HEALTH
	checkEndpointHandler := server.NewhttpHandler()
	checkEndpointHandler.AddInterceptorFunc(http.MethodPost, handler.ErrorHandler)
	checkEndpointHandler.AddHandlerFunc(http.MethodPost, controller.Check)

	http.HandleFunc(CheckEndpoint, checkEndpointHandler.Handle)

	// USERS
	usersHandler := server.NewhttpHandler()
	//Commit only with response code = 200, with other status code: Rollback
	usersHandler.AddInterceptorFunc(http.MethodPost, handler.TransactionScopeHandler)
	usersHandler.AddHandlerFunc(http.MethodPost, controller.CreateUser)
	http.HandleFunc(UsersEndpoint, usersHandler.Handle)

	//LOGIN
	authHandler := server.NewhttpHandler()
	authHandler.AddHandlerFunc(http.MethodPost, controller.Login)
	http.HandleFunc(LoginEndpoint, authHandler.Handle)

	//MESSAGES
	msgsHandler := server.NewhttpHandler()
	msgsHandler.AddInterceptorFunc(http.MethodPost, handler.ValidateTokenHandler)
	//Commit only with response code = 200, with other status code: Rollback
	msgsHandler.AddInterceptorFunc(http.MethodPost, handler.TransactionScopeHandler)
	msgsHandler.AddHandlerFunc(http.MethodPost, controller.SendMessage)

	msgsHandler.AddInterceptorFunc(http.MethodGet, handler.ValidateTokenHandler)
	msgsHandler.AddHandlerFunc(http.MethodGet, controller.GetMessages)
	http.HandleFunc(MessagesEndpoint, msgsHandler.Handle)

	go func() {
		defer FinishServer() //finish components
		defer wg.Done() // let main know we are done cleaning up

		log.Println("Server started at port " + ServerPort)

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatal(err)
		}
	}()

	return srv
}
