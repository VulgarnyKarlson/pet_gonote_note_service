package http

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Addr        string `yaml:"Addr"`
	ReadTimeout int    `yaml:"ReadTimeout"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Get("Adapters.HTTP").Populate(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to http get config: %w", err)
	}
	return &cfg, err
}
