package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kunal697/filesharingcli/internal/handlers"
)

func FileRoute(router *gin.Engine) {
	// Route to upload a file
	router.POST("/upload/:site_name", handlers.Uploadfile)
	// Route to download a file
	router.GET("/getfile/:id", handlers.Getfile)
	// router.GET("/download/:id", handlers.DownloadFile)
	// // Route to delete a file
	// router.DELETE("/file/:id", handlers.DeleteFile)
}
