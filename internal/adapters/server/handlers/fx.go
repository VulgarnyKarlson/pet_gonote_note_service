package handlers

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"adapters/http/handlers",
		fx.Provide(
			New,
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "httpHandlers").Logger()
			return &lg
		}),
	)
}
