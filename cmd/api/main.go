package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bluznierca1/go-chi-swagger-testing/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// init our .env file
	initializeDotEnv()

	// Init our router
	router := router.SetupRouter()

	// define server configs
	srv := &http.Server{
		Addr:    ":9200",
		Handler: router,
	}

	// Run server in goroutine
	go func() {
		// Let's start the server in new goroutine to allow it to listen to termination
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crash: %s \n", err)
		}
	}()

	gracefulShutdown(srv)

}

func initializeDotEnv() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Could not initialize .env file %v", err)
	}
}

// gracefulShutdown makes sure that for 5 seconds server will maintain ongoing processes
// but won't accept new requests (500 response for these)
func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Exiting server gracefully...")

}
