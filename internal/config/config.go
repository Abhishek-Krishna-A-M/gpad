package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/storage"
)

type Config struct {
	GitEnabled bool     `json:"git_enabled"`
	RepoURL    string   `json:"repo_url"`
	Editor     string   `json:"editor"`
	AutoPush   bool     `json:"autopush"`
	Pinned     []string `json:"pinned"`
}

var configCache *Config

func path() string {
	return filepath.Join(storage.GpadDir(), "config.json")
}

func Load() (*Config, error) {
	if configCache != nil {
		return configCache, nil
	}
	data, err := os.ReadFile(path())
	if err != nil {
		c := &Config{
			GitEnabled: false,
			RepoURL:    "",
			Editor:     "",
			AutoPush:   true,
			Pinned:     []string{},
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
	configCache = cfg
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path(), b, 0644)
}

func IsPinned(relPath string) bool {
	cfg, _ := Load()
	for _, p := range cfg.Pinned {
		if p == relPath {
			return true
		}
	}
	return false
}

func Pin(relPath string) error {
	cfg, _ := Load()
	if IsPinned(relPath) {
		return nil
	}
	cfg.Pinned = append(cfg.Pinned, relPath)
	return Save(cfg)
}

func Unpin(relPath string) error {
	cfg, _ := Load()
	updated := cfg.Pinned[:0]
	for _, p := range cfg.Pinned {
		if p != relPath {
			updated = append(updated, p)
		}
	}
	cfg.Pinned = updated
	return Save(cfg)
}
