package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kunal697/filesharingcli/internal/handlers"
)

func SiteRoute(router *gin.Engine) {
	// Route to create a new site
	router.POST("/createsite", handlers.CreateSite)
	router.GET("/site/:site_name", handlers.GetSite)
	router.GET("/sites", handlers.Getallsites)
	router.DELETE("/site/:id", func(c *gin.Context) {
		// Delete a site
	})
}
