package logger

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Level string `yaml:"Level"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Get("Common.Logger").Populate(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to auth get config: %w", err)
	}
	return &cfg, nil
}
