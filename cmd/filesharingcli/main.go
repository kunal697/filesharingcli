package main

import (
	"log"
	"net/http"
	"os" // Import os for accessing environment variables

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kunal697/filesharingcli/internal/db"
	"github.com/kunal697/filesharingcli/internal/routes"
)

func setupRouter() *gin.Engine {
    // Attempt to load the .env file in development environments
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found")
    }

    db.ConnectDB()
    router := gin.Default()

    routes.SiteRoute(router)
    routes.FileRoute(router)

    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello from Vercel!"})
    })

    return router
}


// Entry point for Render deployment
func main() {
	router := setupRouter()

	// Use the PORT environment variable, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback to 8080 if PORT isn't set
	}

	// Run the Gin server on the dynamic port
	log.Printf("Server running on port %s", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
