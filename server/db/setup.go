package db

import (
	"context"
	"fmt"
	"os"
	"picourl-backend/db/generated"
	"picourl-backend/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultPostgresHost = "localhost"
	defaultPostgresPort = "5432"
)

var Queries *generated.Queries

func SetupDb() func() {
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		logger.Log.Info(
			"Postgres environment variable not set, using default",
			"var", "POSTGRES_HOST",
			"defaultValue", defaultPostgresHost,
		)
		postgresHost = defaultPostgresHost
	}

	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresPort == "" {
		logger.Log.Info(
			"Postgres environment variable not set, using default",
			"var", "POSTGRES_PORT",
			"defaultValue", defaultPostgresPort,
		)
		postgresPort = defaultPostgresPort
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		logger.Log.Error("Required postgres environment variable not set", "var", "POSTGRES_USER")
		os.Exit(1)
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		logger.Log.Error("Required postgres environment variable not set", "var", "POSTGRES_PASSWORD")
		os.Exit(1)
	}

	postgresDB := os.Getenv("POSTGRES_DB")
	if postgresDB == "" {
		logger.Log.Error("Required postgres environment variable not set", "var", "POSTGRES_DB")
		os.Exit(1)
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", postgresHost, postgresUser, postgresPassword, postgresDB, postgresPort)

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
