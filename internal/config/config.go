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
	GithubRateLimit int
	Port            int
	PollingInterval int
	MigrateMode     string
	EnableMirror    bool
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{}

	config.GiteaAPIURL = getEnv("GITEA_API_URL", "")
	config.GiteaUser = getEnv("GITEA_USER", "")
	config.GiteaToken = getEnv("GITEA_TOKEN", "")
	config.GithubUser = getEnv("GITHUB_USER", "")
	config.GithubToken = getEnv("GITHUB_TOKEN", "")

	config.GithubRateLimit = getEnvAsInt("GH_RATE_LIMIT", 4990)
	config.Port = getEnvAsInt("PORT", 8080)
	config.PollingInterval = getEnvAsInt("POLLING_INTERVAL_MINUTES", 60)
	config.MigrateMode = getEnv("MIGRATE_MODE", "poll")
	config.EnableMirror = getEnvAsBool("ENABLE_MIRROR", true)

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
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
	if c.GithubRateLimit <= 0 {
		return fmt.Errorf("GH_RATE_LIMIT must be greater than 0")
	}
	return nil
}
