package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"picourl-backend/db/generated"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *generated.Queries

func SetupDb() func() {
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresPort == "" {
		postgresPort = "5432"
	}

	postgresUser, postgresPassword, postgresDB := os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB")

	if postgresUser == "" || postgresPassword == "" || postgresDB == "" {
		log.Fatal("Required database environment variables not set")
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
	m, err := migrate.New("file://db/migrations", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	ctx := context.Background()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", postgresHost, postgresUser, postgresPassword, postgresDB, postgresPort)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	Queries = generated.New(pool)

	return func() {
		pool.Close()
		m.Close()
	}
}
