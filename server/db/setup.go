package db

import (
	"context"
	"net"
	"net/url"
	"os"
	"picourl-backend/config"
	"picourl-backend/db/generated"
	"picourl-backend/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	m, err := migrate.New("file://db/migrations", connectionString)
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
