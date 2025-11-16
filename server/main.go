package main

import (
	"fmt"
	"log"
	"net/http"
	"picourl-backend/db"
	"picourl-backend/handlers"
	"picourl-backend/middleware"
	"picourl-backend/redis"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func handleRateLimited(c *gin.Context, info ratelimit.Info) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"error": fmt.Sprintf("Too many requests. Try again in %s", time.Until(info.ResetTime).String()),
	})
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

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

	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.RedisClient,
		Rate:        5 * time.Minute,
		Limit:       60,
	})

	rateLimitMw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: handleRateLimited,
		KeyFunc:      keyFunc,
	})

	r.Use(rateLimitMw)
	r.Use(middleware.SecurityHeaders)

	r.GET("/:id", handlers.Resolve)
	r.GET("/stats/:id", handlers.Stats)
	r.POST("/shorten", handlers.Shorten)

	r.Run()
}
