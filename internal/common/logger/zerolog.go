package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(conf *Config) {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(conf.Level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)
}

func Get() zerolog.Logger {
	return log.Logger
}
