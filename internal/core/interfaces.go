package core

import (
	"context"

	"gitea-migrate/pkg/models"
)

type GiteaService interface {
	CreateRepo(context context.Context, repo *models.Repository) error
	RepoExists(context context.Context, repoName string) (bool, error)
}

type GitHubService interface {
	FetchRepos(context context.Context) ([]*models.Repository, error)
}

type Poller interface {
	Start()
	Stop()
	AddMirroredRepo(repoName string)
	GetMirroredReposCount() int
}
