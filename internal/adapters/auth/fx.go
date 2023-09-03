package auth

import (
	"context"

	"github.com/rs/zerolog"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"adapters/auth",
		fx.Provide(
			NewAuthConfig,
			NewWrapper,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, wrapper Client) {
				lc.Append(
					fx.Hook{
						OnStart: func(context.Context) error {
							return wrapper.Connect()
						},
						OnStop: func(context.Context) error {
							return wrapper.Close()
						},
					},
				)
			},
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "auth").Logger()
			return &lg
		}),
	)
}
