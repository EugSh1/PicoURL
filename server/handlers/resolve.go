package handlers

import (
	"context"
	"log"
	"net/http"
	"picourl-backend/db"
	"picourl-backend/db/generated"
	"picourl-backend/redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	updateStatsTimeout        = 10 * time.Second
	cacheEligibilityThreshold = 15
	cacheTTL                  = 3 * 24 * time.Hour
)

func updateStats(id, foundLinkUrl, referrer string) {
	ctx, cancel := context.WithTimeout(context.Background(), updateStatsTimeout)
	defer cancel()

	err := db.Queries.CreateClick(ctx, generated.CreateClickParams{
		LinkID:   id,
		Referrer: pgtype.Text{String: referrer, Valid: referrer != ""},
	})

	if err != nil {
		log.Println("Failed to create click for link with id:", id, "error:", err)
		return
	}

	count, err := db.Queries.CountRecentClicks(ctx, id)

	if err != nil {
		log.Println("Failed to count recent clicks for link with id:", id, "error:", err)
		return
	}

	if count >= cacheEligibilityThreshold {
		err := redis.RedisClient.Set(ctx, id, foundLinkUrl, cacheTTL).Err()
		if err != nil {
			log.Println("Failed to add link to the redis cache:", id, "error:", err)
		}
	}
}

func Resolve(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	var foundLink string

	foundLinkFromRedis, err := redis.RedisClient.Get(ctx, id).Result()

	if err == nil {
		foundLink = foundLinkFromRedis
	} else {
		foundLinkFromDb, err := db.Queries.GetLinkById(ctx, id)
		if err != nil {
			log.Println("an error occurred in resolve handler", err)

			var status int
			var errorMessage string

			if err == pgx.ErrNoRows {
				status = http.StatusNotFound
				errorMessage = "Not Found"
			} else {
				status = http.StatusInternalServerError
				errorMessage = "Internal Server Error"
			}

			c.JSON(status, gin.H{
				"error": errorMessage,
			})
			return
		}

		foundLink = foundLinkFromDb.Url
	}

	c.Redirect(http.StatusFound, foundLink)

	referrer := c.Request.Referer()

	go updateStats(id, foundLink, referrer)
}
