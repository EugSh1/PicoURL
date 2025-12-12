package redis

import (
	"context"
	"fmt"
	"os"
	"picourl-backend/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultRedisHost = "localhost"
	defaultRedisPort = "6379"
)

var RedisClient *redis.Client

func SetupRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		logger.Log.Info(
			"Redis environment variable not set, using default",
			"var", "REDIS_HOST",
			"defaultValue", defaultRedisHost,
		)
		redisHost = defaultRedisHost
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		logger.Log.Info(
			"Redis environment variable not set, using default",
			"var", "REDIS_PORT",
			"defaultValue", defaultRedisPort,
		)
		redisPort = defaultRedisPort
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		logger.Log.Error("Required redis environment variable not set", "var", "REDIS_PASSWORD")
		os.Exit(1)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := RedisClient.Ping(ctx).Err()
	if err != nil {
		logger.Log.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
}
