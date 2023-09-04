package circuitbreaker

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"common/circuitbreaker",
		fx.Provide(
			NewConfig,
			NewCircuitBreaker,
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "circuitbreaker").Logger()
			return &lg
		}),
	)
}
