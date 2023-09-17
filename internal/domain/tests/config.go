package tests

import (
	"github.com/ory/dockertest/v3"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/adapters/postgres"
)

type configStorage struct {
	postgres *postgres.Config
}

func newConfigStorage() *configStorage {
	return &configStorage{
		postgres: &postgres.Config{
			UserName: "postgres",
			Password: "1234",
			DBName:   "note_service",
			PoolSize: 10,
		},
	}
}

func (c *configStorage) updatePGConfig(resource *dockertest.Resource) *postgres.Config {
	host, port := getHostPort(resource, "5432/tcp")
	c.postgres.Host = host
	c.postgres.Port = port
	return c.postgres
}

func (c *configStorage) GetPGConfig() *postgres.Config {
	return c.postgres
}
