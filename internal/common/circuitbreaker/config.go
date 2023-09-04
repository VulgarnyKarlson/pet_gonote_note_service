package circuitbreaker

import (
	"fmt"
	"time"

	"go.uber.org/config"
)

type Config struct {
	RecordLength     int           `yaml:"RecordLength"`
	Timeout          time.Duration `yaml:"Timeout"`
	Percentile       float64       `yaml:"Percentile"`
	RecoveryRequests int           `yaml:"RecoveryRequests"`
}

func NewConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Get("Common.CircuitBreaker").Populate(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get CircuitBreaker config: %w", err)
	}
	cfg.Timeout *= time.Millisecond
	return &cfg, nil
}
