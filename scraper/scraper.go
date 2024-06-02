package scraper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
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

// getProxyURL retrieves the proxy URL from the proxy pool URL specified in the environment variable.
func getProxyURL() (string, error) {
	proxyPoolURL := os.Getenv("PROXY_POOL_URL")
	if proxyPoolURL == "" {
		return "", errors.New("PROXY_POOL_URL is not set in .env file")
	}

	resp, err := http.Get(proxyPoolURL)
	if err != nil {
		return "", fmt.Errorf("failed to get proxy from pool: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("proxy pool server returned non-200 status: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read proxy response body: %v", err)
	}

	proxyURL := strings.TrimSpace(string(body))
	if !strings.HasPrefix(proxyURL, "http://") && !strings.HasPrefix(proxyURL, "https://") && !strings.HasPrefix(proxyURL, "socks5://") {
		proxyURL = "http://" + proxyURL
	}

	_, err = url.Parse(proxyURL)
	if err != nil {
		return "", fmt.Errorf("invalid proxy URL: %v", err)
	}

	return proxyURL, nil
}

// Scrape fetches the HTML content from the given URL and saves it to a file.
func Scrape(url string, filename string) error {
	if url == "" {
		return errors.New("URL is required")
	}

	// Initialize Colly collector
	collector := colly.NewCollector()

	// Set proxy if provided via environment variable
	proxyURL, err := getProxyURL()
	if err != nil {
		return err
	}
	if proxyURL != "" {
		// Using the proxy.RoundRobinProxySwitcher to set proxy
		rp, err := proxy.RoundRobinProxySwitcher(proxyURL)
		if err != nil {
			return fmt.Errorf("failed to set proxy: %v", err)
		}
		collector.SetProxyFunc(rp)
	}

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
