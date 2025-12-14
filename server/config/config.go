package config

import (
	"log/slog"
	"os"
	"picourl-backend/logger"
)

const (
	defaultPostgresHost = "localhost"
	defaultPostgresPort = "5432"
	defaultRedisHost    = "localhost"
	defaultRedisPort    = "6379"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
}

func Init() *Config {
	cfg := &Config{
		PostgresHost:     getEnvDefault("POSTGRES_HOST", defaultPostgresHost),
		PostgresPort:     getEnvDefault("POSTGRES_PORT", defaultPostgresPort),
		PostgresUser:     getEnvRequired("POSTGRES_USER"),
		PostgresPassword: getEnvRequired("POSTGRES_PASSWORD"),
		PostgresDB:       getEnvRequired("POSTGRES_DB"),
		RedisHost:        getEnvDefault("REDIS_HOST", defaultRedisHost),
		RedisPort:        getEnvDefault("REDIS_PORT", defaultRedisPort),
		RedisPassword:    getEnvRequired("REDIS_PASSWORD"),
	}

	logger.Log.Info("Config loaded", slog.Group("config",
		slog.String("PostgresHost", cfg.PostgresHost),
		slog.String("PostgresPort", cfg.PostgresPort),
		slog.String("PostgresUser", cfg.PostgresUser),
		slog.String("PostgresPassword", "[REDACTED]"),
		slog.String("PostgresDB", cfg.PostgresDB),
		slog.String("RedisHost", cfg.RedisHost),
		slog.String("RedisPort", cfg.RedisPort),
		slog.String("RedisPassword", "[REDACTED]"),
	))

	return cfg
}

func getEnvDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Log.Info(
			"Environment variable not set, using default",
			"var", key,
			"defaultValue", def,
		)
		return def
	}

	return val
}

func getEnvRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Log.Error("Required environment variable not set", "var", key)
		os.Exit(1)
	}

	return val
}
