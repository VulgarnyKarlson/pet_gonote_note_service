package postgres

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"adapter/postgres",
		fx.Provide(
			NewConfig,
			New,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, wrapper *Pool) {
				lc.Append(
					fx.Hook{
						OnStart: func(context.Context) error {
							return nil
						},
						OnStop: func(context.Context) error {
							wrapper.Close()
							return nil
						},
					},
				)
			},
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "postgres").Logger()
			return &lg
		}),
	)
}
