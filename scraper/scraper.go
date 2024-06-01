package scraper

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.%d Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%d.0 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%d.0) Gecko/20100101 Firefox/%d.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.%d Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%d.0 Mobile/15E148 Safari/604.1",
}

// RandomUserAgent generates a random User-Agent string from a predefined list.
func RandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	ua := userAgents[rand.Intn(len(userAgents))]
	return fmt.Sprintf(ua, rand.Intn(20)+60, rand.Intn(1000)+100, rand.Intn(1000))
}

// Scrape fetches the HTML content from the given URL and saves it to a file.
func Scrape(url string, filename string) error {
	if url == "" {
		return errors.New("URL is required")
	}

	// Initialize Colly collector
	collector := colly.NewCollector()

	// Set up a callback function to be executed when a request is made
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomUserAgent())
	})

	// Create a file to save the HTML content
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set up a callback function to be executed when a page is visited
	collector.OnResponse(func(response *colly.Response) {
		_, err := file.Write(response.Body)
		if err != nil {
			return
		}
	})

	// Visit the URL
	err = collector.Visit(url)
	if err != nil {
		return err
	}

	return nil
}
