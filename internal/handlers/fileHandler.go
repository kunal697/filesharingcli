package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kunal697/filesharingcli/internal/db"
	"github.com/kunal697/filesharingcli/internal/models"
	utils "github.com/kunal697/filesharingcli/internal/utilis"
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

func Uploadfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Validate the token
	claims, err := utils.ValidateToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Verify the site name from the token
	siteName := claims.SiteName

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file", "details": err.Error()})
		return
	}

	// Get the repository details and token from environment variables
	repoOwner := "kunal697"
	repoName := "cliFilesharing"
	accessToken := os.Getenv("GITHUB_TOKEN")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "GitHub token not provided"})
		return
	}
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	file.Filename = file.Filename + "_" + timestampStr

	filePath := fmt.Sprintf("uploads/%s/%s", siteName, file.Filename)

	// Open the file to read its content
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "details": err.Error()})
		return
	}
	defer fileContent.Close()

	// Convert the file content to a byte slice
	fileBytes := make([]byte, file.Size)
	_, err = fileContent.Read(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the file content", "details": err.Error()})
		return
	}

	// Base64 encode the file content (required by GitHub API)
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	// Prepare the request body for uploading the file
	reqBody := map[string]interface{}{
		"message": fmt.Sprintf("Upload file: %s", file.Filename),
		"content": encodedContent,
	}

	// Marshal the request body
	reqBodyBytes, _ := json.Marshal(reqBody)

	// Construct the API URL for uploading the file
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, filePath)

	// Make the PUT request to upload the file
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request", "details": err.Error()})
		return
	}

	// Set headers and send the request
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Check if the upload was successful
	if resp.StatusCode != http.StatusCreated {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to upload file to GitHub",
			"status": resp.StatusCode,
		})
		return
	}

	var site models.Site
	err = db.DB.Where("site_name = ?", siteName).First(&site).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "site not found",
			"status": err.Error(),
		})
		return
	}

	// Save the file metadata to the database
	fileRecord := models.File{
		SiteName: siteName,
		FileName: file.Filename,
		FileURL:  fmt.Sprintf("https://github.com/%s/%s/blob/main/%s", repoOwner, repoName, filePath),
	}

	if err := db.DB.Create(&fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Error in Updating file in database ",
			"status": err.Error(),
		})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully to GitHub",
		"file":    file.Filename,
		"repo":    repoName,
		"path":    filePath,
	})
}

func Getfile(c *gin.Context) {
	fileId := c.Param("id")
	authHeader := c.GetHeader("Authorization")
	claims, err := utils.ValidateToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	var file models.File
	err = db.DB.Where("id = ?", fileId).First(&file).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	if claims.SiteName != file.SiteName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access to file"})
		return
	}

	// Get GitHub token
	accessToken := os.Getenv("GITHUB_TOKEN")
	if accessToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub token not set"})
		return
	}

	// Construct raw GitHub URL
	rawURL := fmt.Sprintf("https://%s@raw.githubusercontent.com/kunal697/cliFilesharing/main/uploads/%s/%s",
		accessToken,
		file.SiteName,
		file.FileName,
	)

	// Make the request
	resp, err := http.Get(rawURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch file"})
		return
	}
	defer resp.Body.Close()

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{
			"error":   "Failed to fetch file",
			"details": string(body),
		})
		return
	}

	// Read file content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File fetched successfully",
		"file":    string(content),
	})
}
