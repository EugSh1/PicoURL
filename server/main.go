package main

import (
	"os"
	"picourl-backend/config"
	"picourl-backend/db"
	"picourl-backend/handlers"
	"picourl-backend/logger"
	"picourl-backend/middleware"
	"picourl-backend/redis"
	"picourl-backend/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	logger.Init()
	logger.Log.Info("Starting PicoURL server...")

	err := godotenv.Load()
	if err != nil {
		logger.Log.Info("Env file not found, using system environment variables")
	}

	cfg := config.Init()

	closeConn := db.SetupDb(cfg)
	defer closeConn()

	redis.SetupRedis(cfg)

	if redis.RedisClient == nil {
		logger.Log.Error("Failed to initialize Redis client, redisClient is nil")
		os.Exit(1)
	}

	go worker.Setup()

	r := gin.New()

	r.Use(sloggin.New(logger.Log))
	r.Use(gin.Recovery())
	r.Use(middleware.CreateRateLimit())
	r.Use(middleware.SecurityHeaders)

	r.GET("/:id", handlers.Resolve)
	r.GET("/stats/:id", handlers.Stats)
	r.POST("/shorten", handlers.Shorten)

	r.Run()
}
