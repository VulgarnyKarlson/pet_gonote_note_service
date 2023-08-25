package config

import (
	"flag"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	path, err := getFilePath()
	if err != nil {
		return nil, errors.Wrap(err, "can not get config path")
	}

	if path == "" {
		return nil, errors.New("can not load config from empty path")
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "can not unmarshal config from file to struct")
	}
	logger.SetupLogger(config.Common.Logger)
	return &config, nil
}

func getFilePath() (string, error) {
	flag.String("config", "./config/prod_config.yaml", "config file path")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return "", err
	}
	viper.AutomaticEnv()

	return viper.GetString("config"), nil
}
