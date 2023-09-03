package rabbitmq

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Host      string `yaml:"Host"`
	Port      int    `yaml:"Port"`
	UserName  string `yaml:"UserName"`
	Password  string `yaml:"Password"`
	QueueName string `yaml:"QueueName"`
}

func NewConfig(cfg config.Provider) (*Config, error) {
	var c Config
	err := cfg.Get("Adapters.RabbitMQ").Populate(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rabbitmq config: %w", err)
	}
	return &c, err
}
