package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitea-migrate/internal/config"
	"gitea-migrate/internal/core"

	"github.com/joho/godotenv"
)

func loadEnv(envFile string) error {
	err := godotenv.Load(envFile)
	if err != nil {
		return fmt.Errorf("Error loading .env file: %w", err)
	}

	requiredEnvVars := []string{"GITEA_API_URL", "GITEA_USER", "GITEA_TOKEN", "GITHUB_USER", "GITHUB_TOKEN"}
	for _, envVar := range requiredEnvVars {
		if value := os.Getenv(envVar); value == "" {
			return fmt.Errorf("Required environment variable %s is not set", envVar)
		}
	}

	return nil
}

type WebhookHandler struct {
	config *config.Config
	poller *core.GithubPoller
}

func NewWebhookHandler(config *config.Config) *WebhookHandler {
	interval := config.PollingInterval
	return &WebhookHandler{
		config: config,
		poller: core.NewGithubPoller(time.Duration(interval), config),
	}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleMigrateWebhook(w, r, h.config)
}

func InitRouter(config *config.Config) (router *http.ServeMux, err error) {
	err = loadEnv(".env")
	if err != nil {
		return nil, fmt.Errorf("Error finding environment variables: %w", err)
	}

	router = http.NewServeMux()

	if config.MigrateMode == "webhook" || config.MigrateMode == "both" {
		handler := NewWebhookHandler(config)
		router.Handle("/migrate-webhook", handler)

		router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, "OK")
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
		})
	}

	return router, nil
}
