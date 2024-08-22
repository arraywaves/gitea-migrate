package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gitea-migrate/api"
	"gitea-migrate/logic"
)

func main() {
	router, err := api.InitRouter()
	if err != nil {
		log.Fatalf("Error initialising router: %v", err)
	}

	mode := os.Getenv("MIRROR_MODE")
	if mode == "" {
		mode = "poll"
	}

	var poller *logic.GithubPoller
	if mode == "poll" || mode == "both" {
		intervalMinutes := 60
		if envInterval := os.Getenv("POLLING_INTERVAL_MINUTES"); envInterval != "" {
			if i, err := strconv.Atoi(envInterval); err == nil {
				intervalMinutes = i
			} else {
				log.Printf("Invalid POLLING_INTERVAL_MINUTES, using default: %v", err)
			}
		}

		pollingInterval := time.Duration(intervalMinutes) * time.Minute
		log.Printf("Setting polling interval to %v", pollingInterval)

		poller = logic.NewGithubPoller(pollingInterval)
		log.Printf("Initial mirrored repos count: %d", poller.GetMirroredReposCount())
		poller.Start()
		defer poller.Stop()
	}

	if mode == "webhook" || mode == "both" {
		log.Println("Webhook endpoint active at /migrate-webhook")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting server on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-waitForInterrupt()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if poller != nil {
		poller.Stop()
	}
	log.Println("Server stopped")
}

func waitForInterrupt() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	return c
}
