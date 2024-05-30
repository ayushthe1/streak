package main

import (
	"log"
	"sync"

	"github.com/ayushthe1/streak/database"
	"github.com/ayushthe1/streak/httpserver"
	"github.com/ayushthe1/streak/ws"
)

func main() {
	database.Connect()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		// Decrement the counter when this goroutine completes
		defer wg.Done()
		ws.StartWebSocketServer()
	}()

	go func() {
		// Decrement the counter when this goroutine completes
		defer wg.Done()
		httpserver.StartHttpServer()
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	log.Println("Both servers have been started and are running concurrently.")
}
