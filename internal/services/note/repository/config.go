package repository

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	CreateNotesBatchSize int `yaml:"CreateNotesBatchSize"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Get("Services.Note.Repository").Populate(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get note repository config: %w", err)
	}
	return &cfg, nil
}
