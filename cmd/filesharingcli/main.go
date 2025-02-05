package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kunal697/filesharingcli/internal/db"
	"github.com/kunal697/filesharingcli/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()
	router := gin.New()

	routes.SiteRoute(router)
	routes.FileRoute(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, world!"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
