package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitea-migrate/internal/api"
	"gitea-migrate/internal/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	router, err := api.NewRouter(config)
	if err != nil {
		log.Fatalf("Error creating router: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Create an error channel to capture server errors
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server is running on port %d", config.Port)
		err := router.Start()
		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case <-stop:
		log.Println("Stopping goroutines, shutting down server...")
	case err := <-serverErrors:
		log.Printf("Server error: %v", err)
	}

	// Perform any cleanup or shutdown operations here
	if err := router.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
