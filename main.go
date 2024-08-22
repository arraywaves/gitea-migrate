package main

import (
	"log"
	"net/http"

	"gitea-migrate/api"
)

func main() {
	router, err := api.InitRouter()
	if err != nil {
		log.Fatalf("Error initialising router: %v", err)
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
