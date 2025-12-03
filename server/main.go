package main

import (
	"log"
	"picourl-backend/db"
	"picourl-backend/handlers"
	"picourl-backend/middleware"
	"picourl-backend/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	closeConn := db.SetupDb()
	defer closeConn()

	redis.SetupRedis()

	if redis.RedisClient == nil {
		log.Fatal("Failed to initialize Redis client")
	}

	r := gin.Default()

	r.Use(middleware.CreateRateLimit())
	r.Use(middleware.SecurityHeaders)

	r.GET("/:id", handlers.Resolve)
	r.GET("/stats/:id", handlers.Stats)
	r.POST("/shorten", handlers.Shorten)

	r.Run()
}
