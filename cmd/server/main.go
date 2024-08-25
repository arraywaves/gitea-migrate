package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitea-migrate/internal/api"
	"gitea-migrate/internal/config"
	"gitea-migrate/internal/core"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
		return
	}

	router, err := api.InitRouter(config)
	if err != nil {
		log.Fatalf("Error initialising router: %v", err)
	}

	if config.MigrateMode == "webhook" || config.MigrateMode == "both" {
		log.Println("Webhook endpoint active at /migrate-webhook")
	}

	var poller *core.GithubPoller
	if config.MigrateMode == "poll" || config.MigrateMode == "both" {
		pollingInterval := time.Duration(config.PollingInterval) * time.Minute
		poller = core.NewGithubPoller(pollingInterval, config)
		log.Printf("Initial mirrored repos count: %d", poller.GetMirroredReposCount())
		poller.Start()
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on :%d", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-waitForInterrupt()
	log.Println("Shutting down server...")

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(context); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if poller != nil {
		poller.Stop()
	}

	log.Println("Server stopped.")
}

func waitForInterrupt() chan os.Signal {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	return exit
}
