package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huyixi/go-web2book/scraper"
)

type Request struct {
	URL string `json:"url"`
}

func ScrapeHandler(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	err := scraper.Scrape(req.URL, "scraped_page.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webpage saved successfully as scraped_page.html"})
}
