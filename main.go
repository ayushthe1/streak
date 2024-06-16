package main

import (
	"log"
	"sync"

	"github.com/ayushthe1/streak/database"
	"github.com/ayushthe1/streak/httpserver"
	"github.com/ayushthe1/streak/kafka"
	"github.com/ayushthe1/streak/upload"
	"github.com/ayushthe1/streak/ws"
)

func main() {
	database.Connect()

	var wg sync.WaitGroup
	wg.Add(4)

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

	go func() {
		defer wg.Done()
		kafka.StartKafkaConsumer()
	}()

	go func() {
		defer wg.Done()
		upload.SetupFileUpload()
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	log.Println("Both servers , kafka, fileUploadConsumer have been started and are running concurrently.")
}
