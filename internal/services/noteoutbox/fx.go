package noteoutbox

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"services/noteoutbox",
		fx.Provide(
			NewRepository,
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "noteoutbox").Logger()
			return &lg
		}),
	)
}
