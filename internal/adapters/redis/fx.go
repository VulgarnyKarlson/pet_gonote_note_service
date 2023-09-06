package redis

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"adapter/redis",
		fx.Provide(
			NewConfig,
			New,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, wrapper Client) {
				lc.Append(
					fx.Hook{
						OnStart: func(context.Context) error {
							return wrapper.HealthCheck()
						},
						OnStop: func(context.Context) error {
							return wrapper.Close()
						},
					},
				)
			},
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "redis").Logger()
			return &lg
		}),
	)
}
