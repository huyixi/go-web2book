package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/huyixi/go-web2book/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.HealthCheckHandler)
	r.GET("/status", handler.HealthCheckHandler)
	r.GET("/crawl/html", handler.CrawlHTMLHandler)

	return r
}
