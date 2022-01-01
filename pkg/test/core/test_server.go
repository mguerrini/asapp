package core

import (
	"context"
	"github.com/challenge/pkg/controller"
	"github.com/challenge/pkg/modules/logger"
	middle "github.com/challenge/pkg/modules/server"
	"github.com/challenge/pkg/server"
	"log"
	"net/http"
	"sync"
	"time"
)

func StopTestServer(w *sync.WaitGroup, srv *http.Server) {
	if r := recover(); r != nil {
		logger.Error("Error on finished test" ,r.(error))
	}

	srv.Shutdown(context.TODO())
	srv.Close()
	w.Wait()
}

func StartTestServer () (*sync.WaitGroup, *http.Server) {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	h := http.NewServeMux()
	srv := server.StartHttpServerWith(httpServerExitDone, h)

	//wait 1 second
	time.Sleep(time.Duration(1) * time.Second)

	return httpServerExitDone, srv
}

func StartTestServerThatFailsOnCreateUser () (*sync.WaitGroup, *http.Server) {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	srv := startHttpServerThatFailsOnCreateUser(httpServerExitDone)

	//wait 1 second
	time.Sleep(time.Duration(1) * time.Second)

	return httpServerExitDone, srv
}

func startHttpServerThatFailsOnCreateUser(wg *sync.WaitGroup) *http.Server {
	server.StartServer() //initialize components

	h := http.NewServeMux()
	//create server
	srv := new(http.Server)
	srv.Handler = h
	srv.Addr = ":" + server.ServerPort

	//create handler
	controller := controller.NewController() // controller.Handler{}
	handler := server.NewRequestHandler()

	// Configure endpoints
	// USERS
	usersHandler := middle.NewhttpHandler()
	//Commit only with response code = 200, with other status code: Rollback
	usersHandler.AddInterceptorFunc(http.MethodPost, handler.TransactionScopeHandler)
	usersHandler.AddInterceptorFunc(http.MethodPost, alwaysFailHandler)
	usersHandler.AddHandlerFunc(http.MethodPost, controller.CreateUser)
	h.HandleFunc(server.UsersEndpoint, usersHandler.Handle)

	//LOGIN
	authHandler := middle.NewhttpHandler()
	authHandler.AddHandlerFunc(http.MethodPost, controller.Login)
	h.HandleFunc(server.LoginEndpoint, authHandler.Handle)

	go func() {
		defer server.FinishServer() //finish components
		defer wg.Done() // let main know we are done cleaning up

		log.Println("Server started at port " + server.ServerPort)

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatal(err)
		}
	}()

	return srv
}

func alwaysFailHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, next middle.HttpInterceptor)  {
	next.Handle(ctx, w, r)

	panic("Test error")
}
