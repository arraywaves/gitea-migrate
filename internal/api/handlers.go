package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitea-migrate/internal/config"
	"gitea-migrate/internal/core"
	"gitea-migrate/pkg/models"
)

type GithubWebhookPayload struct {
	Action     string `json:"action"`
	Repository struct {
		Name     string `json:"name"`
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
}

type Handler struct {
	config        *config.Config
	giteaService  core.GiteaService
	githubService core.GitHubService
	poller        core.Poller
}

func NewHandler(config *config.Config, giteaService core.GiteaService, githubService core.GitHubService, poller core.Poller) *Handler {
	return &Handler{
		config:        config,
		giteaService:  giteaService,
		githubService: githubService,
		poller:        poller,
	}
}

func (h *Handler) HandleMigrateWebhook(w http.ResponseWriter, r *http.Request) {
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

	if webhookPayload.Repository.Name == "" || webhookPayload.Repository.CloneURL == "" {
		log.Printf("Error: Invalid repository name or clone URL")
		http.Error(w, "Invalid repository name or clone URL", http.StatusBadRequest)
		return
	}

	if webhookPayload.Action != "created" {
		log.Printf("Ignoring non-creation event: %s", webhookPayload.Action)
		w.WriteHeader(http.StatusOK)
		return
	}

	repo := &models.Repository{
		Name:     webhookPayload.Repository.Name,
		CloneURL: webhookPayload.Repository.CloneURL,
		Private:  true, // Assuming all repos are private, adjust if needed
	}

	err = h.giteaService.CreateRepo(r.Context(), repo)
	if err != nil {
		log.Printf("Error creating Gitea mirror: %v", err)
		http.Error(w, fmt.Sprintf("Error creating Gitea mirror: %v", err), http.StatusInternalServerError)
		return
	}

	h.poller.AddMirroredRepo(repo.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Repository mirrored successfully"))
}

func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
