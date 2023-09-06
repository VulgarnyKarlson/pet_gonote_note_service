package redis

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
	DB   int    `yaml:"DB"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Get("Adapters.Redis").Populate(&cfg)
	if err != nil {
		return nil, fmt.Errorf("can not populate redis config: %w", err)
	}

	return &cfg, nil
}
