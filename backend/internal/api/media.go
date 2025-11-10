package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"example.com/app/internal/models"
	"example.com/app/internal/services"
)

const (
	MaxFileSize   = 10 * 1024 * 1024 // 10MB
	MaxVideoSize  = 100 * 1024 * 1024 // 100MB
	UploadDir     = "./uploads"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var allowedVideoTypes = map[string]bool{
	"video/mp4":  true,
	"video/webm": true,
	"video/ogg":  true,
}

var storageService services.StorageService

func InitStorageService() error {
	var err error
	storageService, err = services.NewStorageService()
	return err
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
	if !allowedImageTypes[req.ContentType] && !allowedVideoTypes[req.ContentType] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid content type. Only images and videos are allowed",
		})
		return
	}

	// Validate file name length
	if len(req.FileName) == 0 || len(req.FileName) > 255 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "File name must be between 1 and 255 characters",
		})
		return
	}

	ctx := context.Background()
	uploadURL, fileURL, key, err := storageService.GeneratePresignedUploadURL(ctx, req.FileName, req.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to generate upload URL",
			Error:   err.Error(),
		})
		return
	}

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

	contentType := file.Header.Get("Content-Type")
	
	// Determine max file size based on content type
	maxSize := MaxFileSize
	if allowedVideoTypes[contentType] {
		maxSize = MaxVideoSize
	}

	// Validate file size
	if file.Size > int64(maxSize) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: fmt.Sprintf("File size exceeds maximum allowed size of %dMB", maxSize/(1024*1024)),
		})
		return
	}

	// Validate content type
	if !allowedImageTypes[contentType] && !allowedVideoTypes[contentType] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid file type. Only images and videos are allowed",
		})
		return
	}

	// Get or generate key
	key := c.Query("key")
	if key == "" {
		ext := filepath.Ext(file.Filename)
		ctx := context.Background()
		_, _, key, err = storageService.GeneratePresignedUploadURL(ctx, file.Filename, contentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Message: "Failed to generate file key",
				Error:   err.Error(),
			})
			return
		}
		// Extract just the key from the generated URL
		if strings.Contains(key, "?key=") {
			parts := strings.Split(key, "?key=")
			if len(parts) > 1 {
				key = parts[1]
			}
		}
		// If key doesn't have extension, use filename extension
		if filepath.Ext(key) == "" && ext != "" {
			key = key + ext
		}
	}

	// Open the uploaded file
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to open file",
			Error:   err.Error(),
		})
		return
	}
	defer fileReader.Close()

	// Upload file using storage service
	ctx := context.Background()
	fileURL, err := storageService.UploadFile(ctx, fileReader, key, contentType, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to upload file",
			Error:   err.Error(),
		})
		return
	}

	// Save media metadata to database
	media := models.Media{
		Key:         key,
		URL:         fileURL,
		FileName:    file.Filename,
		ContentType: contentType,
		Size:        file.Size,
	}

	if err := db.Create(&media).Error; err != nil {
		// Log error but don't fail the request since file is already uploaded
		fmt.Printf("Warning: Failed to save media metadata: %v\n", err)
	}

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
// @Produce image/jpeg,image/png,image/gif,image/webp,video/mp4
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
// @Description Delete an uploaded media file and clean up database record (admin only)
// @Tags media
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.DeleteMediaRequest true "Delete media request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/media [delete]
func DeleteMedia(c *gin.Context) {
	var req models.DeleteMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	key := req.Key
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

	// Delete from storage
	ctx := context.Background()
	if err := storageService.DeleteFile(ctx, key); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete file",
			Error:   err.Error(),
		})
		return
	}

	// Delete from database
	var media models.Media
	if err := db.Where("key = ?", key).First(&media).Error; err == nil {
		db.Delete(&media)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}

// CleanupReplacedMedia deletes old media when it's replaced
func CleanupReplacedMedia(oldURL string) error {
	if oldURL == "" {
		return nil
	}

	// Extract key from URL
	// URL format: http://host/api/v1/media/{key} or http://host/{bucket}/{key}
	parts := strings.Split(oldURL, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid URL format")
	}

	var key string
	// Check if it's a media API URL
	for i, part := range parts {
		if part == "media" && i+1 < len(parts) {
			key = strings.Join(parts[i+1:], "/")
			break
		}
	}

	// If not found, try to extract from bucket URL
	if key == "" && len(parts) >= 2 {
		// Assume last parts are bucket/key
		key = parts[len(parts)-1]
	}

	if key == "" {
		return fmt.Errorf("could not extract key from URL")
	}

	ctx := context.Background()
	if err := storageService.DeleteFile(ctx, key); err != nil {
		return fmt.Errorf("failed to delete old media: %w", err)
	}

	// Delete from database
	var media models.Media
	if err := db.Where("key = ?", key).First(&media).Error; err == nil {
		db.Delete(&media)
	}

	return nil
}
