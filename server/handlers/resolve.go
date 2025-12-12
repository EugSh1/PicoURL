package handlers

import (
	"context"
	"net/http"
	"picourl-backend/db"
	"picourl-backend/db/generated"
	"picourl-backend/logger"
	"picourl-backend/redis"
	requestgroup "picourl-backend/request_group"
	"picourl-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	updateStatsTimeout        = 10 * time.Second
	cacheEligibilityThreshold = 5
)

func updateStats(cacheKey, id, url, referrer string) {
	ctx, cancel := context.WithTimeout(context.Background(), updateStatsTimeout)
	defer cancel()

	err := db.Queries.CreateClick(ctx, generated.CreateClickParams{
		LinkID:   id,
		Referrer: pgtype.Text{String: referrer, Valid: referrer != ""},
	})

	if err != nil {
		logger.Log.Error("Failed to create click for link", "linkId", id, "error", err)
		return
	}

	count, err := db.Queries.CountRecentClicks(ctx, id)

	if err != nil {
		logger.Log.Error("Failed to count recent clicks for link", "linkId", id, "error", err)
		return
	}

	if count >= cacheEligibilityThreshold {
		utils.CacheResult(cacheKey, url, 3*24*time.Hour)
	}
}

func Resolve(c *gin.Context) {
	id := c.Param("id")
	referrer := c.Request.Referer()
	ctx := c.Request.Context()

	cacheKey := "url:" + id

	link, err := redis.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if link == "not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		c.Redirect(http.StatusFound, link)

		go updateStats(cacheKey, id, link, referrer)
		return
	}

	result, err, _ := requestgroup.Group.Do(cacheKey, func() (any, error) {
		dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		link, err := db.Queries.GetLinkById(dbCtx, id)
		if err != nil {
			logger.Log.Error("Error getting link by id", "linkId", id, "error", err)

			if err == pgx.ErrNoRows {
				utils.CacheResult(cacheKey, "not_found", 5*time.Minute)
			}

			return nil, err
		}

		url := link.Url

		updateStats(cacheKey, id, url, referrer)

		return url, nil
	})

	if err != nil {
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

	resultLink, ok := result.(string)
	if !ok {
		logger.Log.Error("singleflight returned non-string result without error in resolve handler", "result", result)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.Redirect(http.StatusFound, resultLink)
}
