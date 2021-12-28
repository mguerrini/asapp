package main

import (
	"github.com/challenge/pkg/modules/server"
	"log"
	"net/http"

	"github.com/challenge/pkg/controller"
	host "github.com/challenge/pkg/server"

)

const (
	ServerPort = "8080"
	CheckEndpoint = "/check"
	UsersEndpoint = "/users"
	LoginEndpoint = "/login"
	MessagesEndpoint = "/messages"
)

func main() {
	host.StartServer() //initialize components
	defer host.FinishServer() //finish components

	//create handler
	controller := controller.NewController() // controller.Handler{}
	handler := host.NewRequestHandler()

	// Configure endpoints
	// HEALTH
	checkEndpointHandler := server.NewhttpHandler()
	checkEndpointHandler.AddMethodHandlerFunc(http.MethodPost, controller.Check)

	http.HandleFunc(CheckEndpoint, checkEndpointHandler.Handle)

	// USERS
	usersHandler := server.NewhttpHandler()
	usersHandler.AddMethodHandlerFunc(http.MethodPost, controller.CreateUser)
	http.HandleFunc(UsersEndpoint, usersHandler.Handle)

	//LOGIN
	authHandler := server.NewhttpHandler()
	authHandler.AddMethodHandlerFunc(http.MethodPost, controller.Login)
	http.HandleFunc(LoginEndpoint, authHandler.Handle)

	//MESSAGES
	msgsHandler := server.NewhttpHandler()
	msgsHandler.AddCancelMethodHandlerFunc(http.MethodPost, handler.ValidateUserHandler)
	msgsHandler.AddMethodHandlerFunc(http.MethodPost, controller.SendMessage)

	msgsHandler.AddCancelMethodHandlerFunc(http.MethodGet, handler.ValidateUserHandler)
	msgsHandler.AddMethodHandlerFunc(http.MethodGet, controller.GetMessages)
	http.HandleFunc(MessagesEndpoint, msgsHandler.Handle)

	// Start server
	log.Println("Server started at port " + ServerPort)
	log.Fatal(http.ListenAndServe(":" + ServerPort, nil))
}
