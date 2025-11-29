package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

type Config struct {
	GitEnabled bool   `json:"git_enabled"`
	RepoURL    string `json:"repo_url"`
	Editor     string `json:"editor"`
	AutoPush   bool   `json:"autopush"`
}

var configCache *Config

func path() string {
	return filepath.Join(storage.GpadDir(), "config.json")
}

func Load() (*Config, error) {
	if configCache != nil {
		return configCache, nil
	}

	file := path()
	data, err := os.ReadFile(file)
	if err != nil {
		c := &Config{
			GitEnabled: false,
			RepoURL:    "",
			Editor:     "",
		}
		configCache = c
		return c, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	configCache = &cfg
	return &cfg, nil
}

func Save(cfg *Config) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path(), b, 0644)
}

