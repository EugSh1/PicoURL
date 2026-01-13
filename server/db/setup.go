package db

import (
	"context"
	"embed"
	"net"
	"net/url"
	"os"
	"picourl-backend/config"
	"picourl-backend/db/generated"
	"picourl-backend/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var Queries *generated.Queries

func SetupDb(cfg *config.Config) func() {
	q := url.Values{}
	q.Set("sslmode", "disable")

	dsnUrl := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.PostgresUser, cfg.PostgresPassword),
		Host:     net.JoinHostPort(cfg.PostgresHost, cfg.PostgresPort),
		Path:     cfg.PostgresDB,
		RawQuery: q.Encode(),
	}
	connectionString := dsnUrl.String()

	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		logger.Log.Error("Failed to create iofs source driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, connectionString)
	if err != nil {
		logger.Log.Error("Failed to create migrate instance", "error", err)
		os.Exit(1)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Log.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		logger.Log.Error("Failed to parse pgxpool config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Log.Error("Failed to create pgxpool", "error", err)
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Log.Error("Failed to connect to Postgres", "error", err)
		os.Exit(1)
	}

	Queries = generated.New(pool)

	return func() {
		pool.Close()
		m.Close()
	}
}
