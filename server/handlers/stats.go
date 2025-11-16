package handlers

import (
	"log"
	"net/http"
	"picourl-backend/db"

	"github.com/gin-gonic/gin"
)

func Stats(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	stats, err := db.Queries.GetWeeklyClickStats(ctx, id)
	if err != nil {
		log.Println("an error occurred in stats handler", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
