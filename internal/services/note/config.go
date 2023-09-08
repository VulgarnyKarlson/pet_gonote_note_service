package note

import (
	"fmt"

	"go.uber.org/config"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note/repository"
)

type Config struct {
	MaxTitleLength   int               `yaml:"MaxTitleLength"`
	MaxContentLength int               `yaml:"MaxContentLength"`
	Repisotory       repository.Config `yaml:"Repository,omitempty"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	cfg := new(Config)
	err := provider.Get("Services.Note").Populate(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load note config: %w", err)
	}
	return cfg, nil
}
