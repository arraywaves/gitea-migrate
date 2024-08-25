package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitea-migrate/internal/config"
	"gitea-migrate/pkg/models"

	"golang.org/x/time/rate"
)

type GitHubServiceImpl struct {
	config      *config.Config
	rateLimiter *RateLimiter
}

func NewGitHubService(config *config.Config) *GitHubServiceImpl {
	rl := NewRateLimiter(rate.Every(time.Hour/time.Duration(config.GithubRateLimit)), config.GithubRateLimit)
	return &GitHubServiceImpl{
		config:      config,
		rateLimiter: rl,
	}
}

func (g *GitHubServiceImpl) FetchRepos(context context.Context) ([]*models.Repository, error) {
	if err := g.rateLimiter.Wait(context); err != nil {
		return nil, fmt.Errorf("Rate limit exceeded: %v", err)
	}

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
