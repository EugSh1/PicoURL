package db

import (
	"context"
	"fmt"
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
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
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

	ctx := context.Background()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.PostgresPort)

	pool, err := pgxpool.New(ctx, dsn)
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
