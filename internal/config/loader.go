package config

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/config"
	"go.uber.org/fx"
)

type ResultConfig struct {
	fx.Out

	Provider config.Provider
	Config   *Config
}

func NewConfig() (ResultConfig, error) {
	path, err := getFilePath()
	if err != nil {
		return ResultConfig{}, fmt.Errorf("can not get config path: %w", err)
	}

	loader, err := config.NewYAML(config.File(path))
	if err != nil {
		return ResultConfig{}, fmt.Errorf("can not create config loader: %w", err)
	}

	var cfg Config
	err = loader.Get("").Populate(&cfg)
	if err != nil {
		return ResultConfig{}, fmt.Errorf("can not populate config: %w", err)
	}

	return ResultConfig{
		Provider: loader,
		Config:   &cfg,
	}, nil
}

func getFilePath() (string, error) {
	flag.String("config", "./config/local_config.yaml", "config file path")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return "", err
	}
	viper.AutomaticEnv()

	return viper.GetString("config"), nil
}
