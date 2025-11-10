package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"example.com/app/internal/models"
)

const (
	MaxFileSize   = 10 * 1024 * 1024 // 10MB
	UploadDir     = "./uploads"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// GetUploadURL godoc
// @Summary Get presigned URL for file upload
// @Description Get a presigned URL for uploading files to object storage (admin only)
// @Tags media
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UploadURLRequest true "Upload URL request"
// @Success 200 {object} models.UploadURLResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/media/upload-url [post]
func GetUploadURL(c *gin.Context) {
	var req models.UploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate content type
	if !allowedImageTypes[req.ContentType] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid content type. Only images are allowed",
		})
		return
	}

	// Generate unique key for the file
	ext := filepath.Ext(req.FileName)
	key := fmt.Sprintf("dishes/%s/%s%s", time.Now().Format("2006/01"), uuid.New().String(), ext)

	// In a production environment, you would generate a presigned URL for S3/MinIO
	// For now, we'll return a local upload endpoint
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	uploadURL := fmt.Sprintf("%s/api/v1/admin/media/upload?key=%s", baseURL, key)
	fileURL := fmt.Sprintf("%s/api/v1/media/%s", baseURL, key)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data: models.UploadURLResponse{
			UploadURL: uploadURL,
			FileURL:   fileURL,
			Key:       key,
		},
	})
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file directly to the server (admin only)
// @Tags media
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param key query string false "File key/path"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/media/upload [post]
func UploadFile(c *gin.Context) {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "File is required",
			Error:   err.Error(),
		})
		return
	}

	// Validate file size
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: fmt.Sprintf("File size exceeds maximum allowed size of %dMB", MaxFileSize/(1024*1024)),
		})
		return
	}

	// Validate content type
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid file type. Only images are allowed",
		})
		return
	}

	// Get or generate key
	key := c.Query("key")
	if key == "" {
		ext := filepath.Ext(file.Filename)
		key = fmt.Sprintf("dishes/%s/%s%s", time.Now().Format("2006/01"), uuid.New().String(), ext)
	}

	// Ensure upload directory exists
	uploadPath := filepath.Join(UploadDir, filepath.Dir(key))
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create upload directory",
			Error:   err.Error(),
		})
		return
	}

	// Save the file
	filePath := filepath.Join(UploadDir, key)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to save file",
			Error:   err.Error(),
		})
		return
	}

	// Generate file URL
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	fileURL := fmt.Sprintf("%s/api/v1/media/%s", baseURL, key)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "File uploaded successfully",
		Data: models.MediaFile{
			URL:         fileURL,
			FileName:    file.Filename,
			ContentType: contentType,
			Size:        file.Size,
		},
	})
}

// GetMedia godoc
// @Summary Get uploaded media file
// @Description Serve uploaded media files
// @Tags media
// @Produce image/jpeg,image/png,image/gif,image/webp
// @Param filepath path string true "File path"
// @Success 200 {file} binary
// @Failure 404 {object} models.ErrorResponse
// @Router /media/{filepath} [get]
func GetMedia(c *gin.Context) {
	// Get the file path from the URL
	filePath := c.Param("filepath")
	
	// Prevent directory traversal attacks
	if strings.Contains(filePath, "..") {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid file path",
		})
		return
	}

	fullPath := filepath.Join(UploadDir, filePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Serve the file
	c.File(fullPath)
}

// DeleteMedia godoc
// @Summary Delete uploaded media file
// @Description Delete an uploaded media file (admin only)
// @Tags media
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param key query string true "File key/path"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/media [delete]
func DeleteMedia(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "File key is required",
		})
		return
	}

	// Prevent directory traversal attacks
	if strings.Contains(key, "..") {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid file key",
		})
		return
	}

	fullPath := filepath.Join(UploadDir, key)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete file",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}
