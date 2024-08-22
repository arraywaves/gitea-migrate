package logic

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type GithubRepo struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
}

type GithubPoller struct {
	mirroredRepos map[string]bool
	mutex         sync.Mutex
	stopChan      chan struct{}
}

func NewGithubPoller() *GithubPoller {
	poller := &GithubPoller{
		mirroredRepos: make(map[string]bool),
		stopChan:      make(chan struct{}),
	}
	poller.loadMirroredRepos()
	return poller
}

func (p *GithubPoller) Start() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.checkForNewRepos()
			case <-p.stopChan:
				return
			}
		}
	}()
}

func (p *GithubPoller) Stop() {
	close(p.stopChan)
}

func (p *GithubPoller) checkForNewRepos() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	url := "https://api.github.com/user/repos"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+githubToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching repos: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var repos []GithubRepo
	json.Unmarshal(body, &repos)

	for _, repo := range repos {
		p.mutex.Lock()
		if !p.mirroredRepos[repo.Name] {
			err := CreateGiteaMirror(repo.Name, repo.CloneURL)
			if err != nil {
				log.Printf("Error mirroring repo %s: %v", repo.Name, err)
			} else {
				p.mirroredRepos[repo.Name] = true
				p.saveMirroredRepos()
				log.Printf("Successfully mirrored repo: %s", repo.Name)
			}
		}
		p.mutex.Unlock()
	}
}

func (p *GithubPoller) loadMirroredRepos() {
	file, err := os.ReadFile("mirrored_repos.json")
	if err != nil {
		log.Printf("No existing mirrored repos file found")
		return
	}
	json.Unmarshal(file, &p.mirroredRepos)
}

func (p *GithubPoller) saveMirroredRepos() {
	file, _ := json.Marshal(p.mirroredRepos)
	os.WriteFile("mirrored_repos.json", file, 0644)
}
