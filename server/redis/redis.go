package redis

import (
	"context"
	"fmt"
	"os"
	"picourl-backend/config"
	"picourl-backend/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func SetupRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := RedisClient.Ping(ctx).Err()
	if err != nil {
		logger.Log.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
}
