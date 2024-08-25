package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitea-migrate/internal/config"
	"gitea-migrate/pkg/models"
)

type GitHubServiceImpl struct {
	config *config.Config
}

func NewGitHubService(config *config.Config) *GitHubServiceImpl {
	return &GitHubServiceImpl{config: config}
}

func (g *GitHubServiceImpl) FetchRepos(context context.Context) ([]*models.Repository, error) {
	githubURL := "https://api.github.com/user/repos"

	req, err := http.NewRequestWithContext(context, "GET", githubURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "token "+g.config.GithubToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching repos: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var repos []*models.Repository
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return repos, nil
}
