package repository

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"services/note/repository",
		fx.Provide(
			NewConfig,
			NewRepository,
		),
		fx.Decorate(func(log zerolog.Logger) zerolog.Logger {
			return log.With().Str("module", "note_repository").Logger()
		}),
	)
}
