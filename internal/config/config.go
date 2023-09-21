package config

import (
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/auth"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/rabbitmq"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/redis"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/server"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/circuitbreaker"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/common/logger"
	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/services/note"
)

type Config struct {
	Services struct {
		Note *note.Config `yaml:"Note"`
	} `yaml:"Services"`
	Adapters struct {
		Auth     *auth.Config     `yaml:"Auth"`
		Server   *server.Config   `yaml:"Server"`
		Postgres *postgres.Config `yaml:"Postgres"`
		RabbitMQ *rabbitmq.Config `yaml:"RabbitMQ"`
		Redis    *redis.Config    `yaml:"Redis"`
	} `yaml:"Adapters"`
	Common struct {
		Logger         *logger.Config         `yaml:"Logger"`
		CircuitBreaker *circuitbreaker.Config `yaml:"CircuitBreaker"`
	} `yaml:"Common"`
}
