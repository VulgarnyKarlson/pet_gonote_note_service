package http

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"adapters/http",
		fx.Provide(
			NewConfig,
			NewServer,
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "http").Logger()
			return &lg
		}),
	)
}
