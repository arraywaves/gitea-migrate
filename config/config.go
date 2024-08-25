package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GiteaAPIURL     string
	GiteaUser       string
	GiteaToken      string
	GithubUser      string
	GithubToken     string
	Port            int
	PollingInterval int
	MigrateMode     string
	EnableMirror    bool
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{}

	config.GiteaAPIURL = os.Getenv("GITEA_API_URL")
	config.GiteaUser = os.Getenv("GITEA_USER")
	config.GiteaToken = os.Getenv("GITEA_TOKEN")
	config.GithubUser = os.Getenv("GITHUB_USER")
	config.GithubToken = os.Getenv("GITHUB_TOKEN")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		config.Port = 8080
	} else {
		config.Port = port
	}

	interval, err := strconv.Atoi(os.Getenv("POLLING_INTERVAL_MINUTES"))
	if err != nil {
		config.PollingInterval = 60
	} else {
		config.PollingInterval = interval
	}

	config.MigrateMode = os.Getenv("MIGRATE_MODE")
	if config.MigrateMode == "" {
		config.MigrateMode = "poll"
	}

	enableMirror, err := strconv.ParseBool(os.Getenv("ENABLE_MIRROR"))
	if err != nil {
		config.EnableMirror = true
	} else {
		config.EnableMirror = enableMirror
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.GiteaAPIURL == "" {
		return fmt.Errorf("GITEA_API_URL is required")
	}
	if c.GiteaUser == "" {
		return fmt.Errorf("GITEA_USER is required")
	}
	if c.GiteaToken == "" {
		return fmt.Errorf("GITEA_TOKEN is required")
	}
	if c.GithubUser == "" {
		return fmt.Errorf("GITHUB_USER is required")
	}
	if c.GithubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}
	return nil
}
