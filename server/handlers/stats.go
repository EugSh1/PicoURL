package handlers

import (
	"context"
	"log"
	"net/http"
	"picourl-backend/db"
	"picourl-backend/redis"
	requestgroup "picourl-backend/request_group"
	"picourl-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorGettingStats struct {
	statusCode int
	message    string
}

func (e ErrorGettingStats) Error() string {
	return e.message
}

func Stats(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	cacheKey := "stats:" + id

	if stats, err := redis.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		if stats == "not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		c.Data(http.StatusOK, "application/json", []byte(stats))
		return
	}

	result, err, _ := requestgroup.Group.Do(cacheKey, func() (any, error) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		stats, err := db.Queries.GetWeeklyClickStats(dbCtx, id)
		if err != nil {
			log.Println("an error occurred in stats handler", err)
			return nil, &ErrorGettingStats{
				message:    "Internal Server Error",
				statusCode: http.StatusInternalServerError,
			}
		}

		if stats == nil {
			utils.CacheResult(cacheKey, "not_found", time.Minute)
			return nil, &ErrorGettingStats{
				message:    "Not Found",
				statusCode: http.StatusNotFound,
			}
		}

		utils.CacheResult(cacheKey, stats, 5*time.Minute)

		return stats, nil
	})

	if err != nil {
		errResult, ok := err.(*ErrorGettingStats)

		if ok {
			c.JSON(errResult.statusCode, gin.H{
				"error": errResult.message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
