package auth

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Address string `yaml:"Address"`
}

func NewAuthConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Get("Adapters.Auth").Populate(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to auth get config: %w", err)
	}
	return &cfg, nil
}
