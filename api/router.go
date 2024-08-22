package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

func InitRouter() (router *http.ServeMux, err error) {
	err = loadEnv(".env")
	if err != nil {
		return nil, fmt.Errorf("Error finding environment variables: %w", err)
	}

	router = http.NewServeMux()

	mode := os.Getenv("MIRROR_MODE")
	if mode == "webhook" || mode == "both" || mode == "" {
		router.HandleFunc("/migrate-webhook", handleMigrateWebhook)
	}

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "OK")
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	return router, nil
}
