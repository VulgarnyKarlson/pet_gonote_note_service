package auth

import (
	"fmt"
	"time"

	"go.uber.org/config"
)

type Config struct {
	Address           string        `yaml:"Address"`
	BackupStorageTime time.Duration `yaml:"BackupStorageTime"`
}

func NewAuthConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Get("Adapters.Auth").Populate(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to auth get config: %w", err)
	}
	cfg.BackupStorageTime *= time.Second
	return &cfg, nil
}
