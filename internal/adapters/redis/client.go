package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type ClientImpl struct {
	client *redis.Client
	logger *zerolog.Logger
	config *Config
}

type Client interface {
	HealthCheck() error
	Close() error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}

func New(cfg *Config, logger *zerolog.Logger) Client {
	dbConn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     dbConn,
		Password: "",
		DB:       cfg.DB,
	})

	return &ClientImpl{
		client: client,
		config: cfg,
		logger: logger,
	}
}

func (c *ClientImpl) HealthCheck() error {
	_, err := c.client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientImpl) Close() error {
	return c.client.Close()
}

func (c *ClientImpl) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *ClientImpl) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}
