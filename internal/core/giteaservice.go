package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitea-migrate/internal/config"
	"gitea-migrate/pkg/models"
)

type GiteaServiceImpl struct {
	config *config.Config
}

func NewGiteaService(config *config.Config) *GiteaServiceImpl {
	return &GiteaServiceImpl{config: config}
}

func (g *GiteaServiceImpl) CreateRepo(context context.Context, repo *models.Repository) error {
	userURL := fmt.Sprintf("%s/user", g.config.GiteaAPIURL)
	req, err := http.NewRequestWithContext(context, "GET", userURL, nil)
	if err != nil {
		return fmt.Errorf("error creating user request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.config.GiteaToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending user request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code for user request: %d, body: %s", resp.StatusCode, string(body))
	}

	var userData struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return fmt.Errorf("error decoding user response: %v", err)
	}

	migrateURL := fmt.Sprintf("%s/repos/migrate", g.config.GiteaAPIURL)
	payload := map[string]interface{}{
		"repo_name":     repo.Name,
		"clone_addr":    repo.CloneURL,
		"mirror":        g.config.EnableMirror,
		"private":       repo.Private,
		"auth_username": g.config.GithubUser,
		"auth_password": g.config.GithubToken,
		"service":       "github",
		"wiki":          true,
		"labels":        true,
		"issues":        true,
		"pull_requests": true,
		"releases":      true,
		"repo_owner":    g.config.GiteaUser,
		"uid":           userData.ID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err = http.NewRequestWithContext(context, "POST", migrateURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.config.GiteaToken))

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending migration request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully created mirrored repository: %s", repo.Name)
	return nil
}
