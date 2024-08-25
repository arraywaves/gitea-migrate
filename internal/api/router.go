package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitea-migrate/internal/config"
	"gitea-migrate/internal/core"
)

type Router struct {
	config        *config.Config
	giteaService  core.GiteaService
	githubService core.GitHubService
	poller        core.Poller
	handler       *Handler
	server        *http.Server
}

func NewRouter(config *config.Config) (*Router, error) {
	giteaService := core.NewGiteaService(config)
	githubService := core.NewGitHubService(config)

	interval := time.Duration(config.PollingInterval) * time.Minute
	poller := core.NewGithubPoller(interval, config, giteaService, githubService)

	handler := NewHandler(config, giteaService, githubService, poller)

	return &Router{
		config:        config,
		giteaService:  giteaService,
		githubService: githubService,
		poller:        poller,
		handler:       handler,
	}, nil
}

func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	if r.config.MigrateMode == "webhook" || r.config.MigrateMode == "both" {
		mux.HandleFunc("/migrate-webhook", r.handler.HandleMigrateWebhook)
		log.Println("Webhook endpoint active at /migrate-webhook")
	}

	mux.HandleFunc("/health", r.handler.HandleHealthCheck)

	return mux
}

func (r *Router) Start() error {
	mux := r.SetupRoutes()

	if r.config.MigrateMode == "poll" || r.config.MigrateMode == "both" {
		r.poller.Start()
		log.Printf("Polling started with interval of %d minutes", r.config.PollingInterval)
	}

	r.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", r.config.Port),
		Handler: mux,
	}

	log.Printf("Starting server on :%d", r.config.Port)
	return r.server.ListenAndServe()
}

func (r *Router) Stop() error {
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := r.server.Shutdown(context); err != nil {
		return fmt.Errorf("Server shutdown failed: %v", err)
	}

	if r.poller != nil {
		r.poller.Stop()
	}

	return nil
}
