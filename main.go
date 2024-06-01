package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/huyixi/go-web2book/handler"
)

func main() {
	r := gin.Default()

	// API 路由
	r.POST("/scrape", handlers.ScrapeHandler)

	fmt.Println("Server started at :8080")
	r.Run(":8080")
}
