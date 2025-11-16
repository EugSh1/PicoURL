package handlers

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"picourl-backend/db"
	"picourl-backend/db/generated"
	"picourl-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewUrlDto struct {
	Url string `json:"url" binding:"required,url"`
}

func Shorten(c *gin.Context) {
	var body NewUrlDto

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	clientBaseUrl := os.Getenv("CLIENT_BASE_URL")

	if clientBaseUrl != "" {
		baseURL, err := url.Parse(clientBaseUrl)
		inputURL, err2 := url.Parse(body.Url)

		if err == nil && err2 == nil {
			if strings.EqualFold(baseURL.Host, inputURL.Host) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "you can’t shorten a link that already belongs to PicoURL.",
				})
				return
			}
		}
	}

	shortUrlId, err := utils.GenerateUniqueShortId(c.Request.Context())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate unique short ID for a url",
		})
		return
	}

	err = db.Queries.CreateLink(c.Request.Context(), generated.CreateLinkParams{
		ID:  shortUrlId,
		Url: body.Url,
	})
	if err != nil {
		log.Println("an error occurred in shorten handler while creating a link", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save URL",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"shortUrlId": shortUrlId,
	})
}
