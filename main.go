package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

type Request struct {
	URL string `json:"url"`
}

func scrapeHandler(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	// Initialize Colly collector
	collector := colly.NewCollector()

	// Create a file to save the HTML content
	file, err := os.Create("scraped_page.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer file.Close()

	// Set up a callback function to be executed when a page is visited
	collector.OnResponse(func(response *colly.Response) {
		_, err := file.Write(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the content"})
			return
		}
	})

	// Visit the URL
	err = collector.Visit(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the URL"})
		return
	}

	// Respond to the client
	c.JSON(http.StatusOK, gin.H{"message": "Webpage saved successfully as scraped_page.html"})
}

func main() {
	r := gin.Default()
	r.POST("/scrape", scrapeHandler)
	fmt.Println("Server started at :8080")
	r.Run(":8080")
}
