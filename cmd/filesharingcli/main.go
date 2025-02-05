package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kunal697/filesharingcli/internal/db"
	"github.com/kunal697/filesharingcli/internal/routes"
)

func setupRouter() *gin.Engine {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	db.ConnectDB()
	router := gin.New()

	routes.SiteRoute(router)
	routes.FileRoute(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
	})

	return router
}

// Vercel expects an exported function named Handler
var Handler = setupRouter()
