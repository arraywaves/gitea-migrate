package main

import (
	"log"
	"net/http"
	"os"

	"gitea-migrate/api"
	"gitea-migrate/logic"
)

func main() {
	router, err := api.InitRouter()
	if err != nil {
		log.Fatalf("Error initialising router: %v", err)
	}

	poller := logic.NewGithubPoller()
	poller.Start()
	defer poller.Stop()

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
}
