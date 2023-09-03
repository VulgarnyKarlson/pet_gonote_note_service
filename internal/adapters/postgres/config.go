package postgres

import (
	"fmt"

	"go.uber.org/config"
)

type Config struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
	PoolSize int    `yaml:"PoolSize"`
	SSLMode  string `yaml:"SSLMode"`
}

func NewConfig(cfg config.Provider) (*Config, error) {
	var c Config
	err := cfg.Get("Adapters.Postgres").Populate(&c)
	if err != nil {
		return nil, fmt.Errorf("can not populate db config: %w", err)
	}
	return &c, nil
}
