package main

import (
	"sync"

	host "github.com/challenge/pkg/server"
)

func main() {
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)

	host.StartHttpServer(httpServerExitDone)

	httpServerExitDone.Wait()
}

