package note

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"services/note",
		fx.Provide(
			NewConfig,
			NewService,
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "note").Logger()
			return &lg
		}),
	)
}
