package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitea-migrate/logic"
)

type GithubWebhookPayload struct {
	Repository struct {
		Name     string `json:"name"`
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
}

func handleMigrateWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var webhookPayload GithubWebhookPayload
	if err := json.Unmarshal(payload, &webhookPayload); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Error parsing webhook payload", http.StatusBadRequest)
		return
	}

	log.Printf("Debug - Received webhook for repository: %s", webhookPayload.Repository.Name)
	log.Printf("Debug - Clone URL: %s", webhookPayload.Repository.CloneURL)

	if webhookPayload.Repository.Name == "" || webhookPayload.Repository.CloneURL == "" {
		log.Printf("Error: Invalid repository name or clone URL")
		http.Error(w, "Invalid repository name or clone URL", http.StatusBadRequest)
		return
	}

	err = logic.CreateGiteaMirror(webhookPayload.Repository.Name, webhookPayload.Repository.CloneURL)
	if err != nil {
		log.Printf("Error creating Gitea mirror: %v", err)
		http.Error(w, fmt.Sprintf("Error creating Gitea mirror: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Repository mirrored successfully"))
}
