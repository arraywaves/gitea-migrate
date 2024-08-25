package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitea-migrate/config"
)

func CreateGiteaRepo(repoName, cloneURL string, config *config.Config) error {
	userURL := fmt.Sprintf("%s/user", config.GiteaAPIURL)
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return fmt.Errorf("Error creating user request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.GiteaToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending user request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Unexpected status code for user request: %d, body: %s", resp.StatusCode, string(body))
	}

	var userData struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return fmt.Errorf("Error decoding user response: %v", err)
	}

	migrateURL := fmt.Sprintf("%s/repos/migrate", config.GiteaAPIURL)
	payload := map[string]interface{}{
		"repo_name":     repoName,
		"clone_addr":    cloneURL,
		"mirror":        config.EnableMirror,
		"private":       true,
		"auth_username": config.GithubUser,
		"auth_password": config.GithubToken,
		"service":       "github",
		"wiki":          true,
		"labels":        true,
		"issues":        true,
		"pull_requests": true,
		"releases":      true,
		"repo_owner":    config.GiteaUser,
		"uid":           userData.ID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %v", err)
	}

	req, err = http.NewRequest("POST", migrateURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.GiteaToken))

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending migration request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully created mirrored repository: %s", repoName)
	return nil
}
