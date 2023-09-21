package server

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"adapters/server",
		fx.Provide(
			NewConfig,
			NewServer,
			mux.NewRouter,
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "server").Logger()
			return &lg
		}),
	)
}
