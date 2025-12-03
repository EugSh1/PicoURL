package middleware

import (
	"fmt"
	"net/http"
	"picourl-backend/redis"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func handleRateLimited(c *gin.Context, info ratelimit.Info) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"error": fmt.Sprintf("Too many requests. Try again in %s", time.Until(info.ResetTime).String()),
	})
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func CreateRateLimit() gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.RedisClient,
		Rate:        5 * time.Minute,
		Limit:       60,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: handleRateLimited,
		KeyFunc:      keyFunc,
	})
}
