package rabbitmq

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"auth/rabbitmq",
		fx.Provide(
			NewConfig,
			NewPublisher,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, wrapper *Publisher) {
				lc.Append(
					fx.Hook{
						OnStart: func(context.Context) error {
							return wrapper.Open()
						},
						OnStop: func(context.Context) error {
							return wrapper.Close()
						},
					},
				)
			},
		),
		fx.Decorate(func(log *zerolog.Logger) *zerolog.Logger {
			lg := log.With().Str("module", "rabbitmq").Logger()
			return &lg
		}),
	)
}
