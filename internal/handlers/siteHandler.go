package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kunal697/filesharingcli/internal/db"
	"github.com/kunal697/filesharingcli/internal/models"
	"github.com/kunal697/filesharingcli/internal/utilis"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
)

func init() {
	// Try to load from .env file, but don't error if it doesn't exist
	godotenv.Load()

	// Set default values or use environment variables
	if os.Getenv("DATABASE_URL") == "" {
		// Use a default or panic
		panic("DATABASE_URL environment variable is required")
	}
	if os.Getenv("GITHUB_TOKEN") == "" {
		panic("GITHUB_TOKEN environment variable is required")
	}
}

func CreateSite(c *gin.Context) {
	var input struct {
		SiteName string `json:"site_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if site name already exists
	var existingSite models.Site
	if err := db.DB.Where("name = ?", input.SiteName).First(&existingSite).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Site name already exists"})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new site with hashed password
	site := models.Site{
		SiteName: input.SiteName,
		Password: string(hashedPassword),
	}

	if err := db.DB.Create(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create site"})
		return
	}
	authToken, err := utils.GenerateToken(site.SiteName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Site created successfully",
		"auth_token": authToken,
	})
}

// GetSite retrieves details of a site
func GetSite(c *gin.Context) {
	siteName := c.Param("site_name")
	password := c.Query("password")

	// Find the site
	var site models.Site
	if err := db.DB.Where("site_name = ?", siteName).First(&site).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check password
	if err := utils.VerifyPassword(site.Password, password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong password"})
		return
	}

	var files []models.File // Fix: Use slice instead of single struct
	err := db.DB.Select("id, site_name, file_name, created_at").Where("site_name = ?", siteName).Find(&files).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching files"})
		return
	}

	authToken, err := utils.GenerateToken(site.SiteName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"site":       site,
		"files":      files, // Fix: Correct key name from "fiels" to "files"
		"auth_token": authToken,
	})
}

func Getallsites(c *gin.Context) {
	var site []models.Site
	if err := db.DB.Find(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sites": site})
}

// DeleteSite deletes a site
func DeleteSite(c *gin.Context) {
	siteName := c.Param("site_name")
	password := c.Query("password")

	// Find the site
	var site models.Site
	if err := db.DB.Where("site_name = ?", siteName).First(&site).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check password
	if site.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Delete the site
	if err := db.DB.Delete(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site deleted successfully"})
}
