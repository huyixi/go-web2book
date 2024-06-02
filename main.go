package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/huyixi/go-web2book/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := routes.SetupRouter()

	fmt.Println("Server started at :8080")
	r.Run(":8080")
}
