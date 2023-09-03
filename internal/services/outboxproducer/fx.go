package outboxproducer

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"services/outboxproducer",
		fx.Provide(
			NewOutBoxProducer,
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "outboxproducer").Logger()
			return &lg
		}),
	)
}
