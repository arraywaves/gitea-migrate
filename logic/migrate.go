package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CreateGiteaRepo(repoName, cloneURL string) error {
	giteaAPIURL := os.Getenv("GITEA_API_URL")
	giteaUser := os.Getenv("GITEA_USER")
	giteaToken := os.Getenv("GITEA_TOKEN")
	githubUser := os.Getenv("GITHUB_USER")
	githubToken := os.Getenv("GITHUB_TOKEN")

	enableMirrorStr := os.Getenv("ENABLE_MIRROR")
	if enableMirrorStr == "" {
		enableMirrorStr = "true"
	}

	enableMirror, err := strconv.ParseBool(enableMirrorStr)
	if err != nil {
		return fmt.Errorf("Error parsing ENABLE_MIRROR: %v", err)
	}

	if giteaAPIURL == "" || giteaToken == "" || giteaUser == "" || githubUser == "" || githubToken == "" {
		return fmt.Errorf("Missing required environment variables")
	}

	// First, get the user ID
	userURL := fmt.Sprintf("%s/user", giteaAPIURL)
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return fmt.Errorf("Error creating user request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", giteaToken))

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

	// Now create the repository
	url := fmt.Sprintf("%s/repos/migrate", giteaAPIURL)
	payload := map[string]interface{}{
		"repo_name":     repoName,
		"clone_addr":    cloneURL,
		"mirror":        enableMirror,
		"private":       true,
		"auth_username": githubUser,
		"auth_password": githubToken,
		"service":       "github",
		"wiki":          true,
		"labels":        true,
		"issues":        true,
		"pull_requests": true,
		"releases":      true,
		"repo_owner":    giteaUser,
		"uid":           userData.ID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %v", err)
	}

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", giteaToken))

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully created mirrored repository: %s", repoName)
	return nil
}
