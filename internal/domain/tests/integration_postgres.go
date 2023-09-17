package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func (t *TestIntegration) PostgresUP() {
	resource, err := t.pool.Run("postgres", "14.3-bullseye", []string{
		"POSTGRES_DB=note_service", "POSTGRES_USER=postgres", "POSTGRES_PASSWORD=1234",
	})
	if err != nil {
		t.l.Fatal().Err(err).Msg("Could not start postgres")
	}
	t.l.Info().Msg("Postgres started")

	if err = t.pool.Retry(func() error {
		t.Configs.updatePGConfig(resource)
		pgConfig := t.Configs.GetPGConfig()
		postgresURL := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			t.Configs.postgres.UserName, pgConfig.Password, pgConfig.Host, pgConfig.Port, pgConfig.DBName,
		)
		cfg, errPG := pgx.ParseConfig(postgresURL)
		if errPG != nil {
			return fmt.Errorf("could not parse postgres url: %w", errPG)
		}
		con := stdlib.OpenDB(*cfg)
		defer con.Close()
		err = con.Ping()
		if err != nil {
			return err
		}
		t.l.Info().Msg("Postgres connected")

		_ = os.Setenv("POSTGRES_URL", postgresURL)
		abs, _ := filepath.Abs("../")
		migrationsPath := strings.Split(abs, "internal")[0] + "db/migrations"

		err = goose.Up(con, migrationsPath)
		if err != nil {
			t.l.Err(err).Msg("Could not migrate postgres")
		}
		return err
	}); err != nil {
		t.l.Fatal().Err(err).Msg("Could not connect to postgres")
	} else {
		t.l.Info().Msg("Postgres migrated")
	}

	t.resources = append(t.resources, resource)
}
