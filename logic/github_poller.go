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
	interval      time.Duration
}

var PollingInterval = 60 * time.Minute

func NewGithubPoller(interval time.Duration) *GithubPoller {
	poller := &GithubPoller{
		mirroredRepos: make(map[string]bool),
		stopChan:      make(chan struct{}),
		interval:      interval,
	}
	poller.loadMirroredRepos()
	return poller
}

func (p *GithubPoller) Start() {
	go func() {
		p.checkForNewRepos()
		ticker := time.NewTicker(p.interval)
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
	log.Println("Checking for new repos...")

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Println("GITHUB_TOKEN is not set. Skipping check.")
		return
	}
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

	log.Printf("Finished checking. Total mirrored repos: %d", len(p.mirroredRepos))
	p.saveMirroredRepos()
}

func (p *GithubPoller) loadMirroredRepos() {
	file, err := os.ReadFile("mirrored_repos.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No existing mirrored repos file found. Starting with an empty list.")
			return
		}
		log.Printf("Error reading mirrored repos file: %v", err)
		return
	}
	err = json.Unmarshal(file, &p.mirroredRepos)
	if err != nil {
		log.Printf("Error parsing mirrored repos file: %v", err)
		// If we can't parse the file, start with an empty map
		p.mirroredRepos = make(map[string]bool)
	}
	log.Printf("Loaded %d mirrored repos from file", len(p.mirroredRepos))
}

func (p *GithubPoller) saveMirroredRepos() {
	file, err := json.Marshal(p.mirroredRepos)
	if err != nil {
		log.Printf("Error marshaling mirrored repos: %v", err)
		return
	}
	err = os.WriteFile("mirrored_repos.json", file, 0644)
	if err != nil {
		log.Printf("Error saving mirrored repos file: %v", err)
	} else {
		log.Printf("Saved %d mirrored repos to file", len(p.mirroredRepos))
	}
}

func (p *GithubPoller) GetMirroredReposCount() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return len(p.mirroredRepos)
}

func (p *GithubPoller) AddMirroredRepo(repoName string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.mirroredRepos[repoName] = true
	p.saveMirroredRepos()
}