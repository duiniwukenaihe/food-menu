package api

import (
	"net/http"
	"strconv"
	"example.com/app/internal/models"
	"github.com/gin-gonic/gin"
)

// GetContent godoc
// @Summary Get all published content
// @Description Retrieve paginated list of published content
// @Tags content
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query int false "Filter by category ID"
// @Param search query string false "Search in title and description"
// @Success 200 {object} models.PaginatedResponse
// @Router /content [get]
func GetContent(c *gin.Context) {
	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid pagination parameters",
			Error:   err.Error(),
		})
		return
	}

	var contents []models.Content
	var total int64

	query := db.Model(&models.Content{}).Where("is_published = ?", true)

	// Filter by category
	if categoryID := c.Query("category"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Search functionality
	if search := c.Query("search"); search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total records
	query.Count(&total)

	// Get paginated results
	if err := query.Preload("Category").Preload("Author").
		Offset(pagination.GetOffset()).
		Limit(pagination.Limit).
		Order("created_at DESC").
		Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch content",
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	var contentResponses []models.ContentResponse
	for _, content := range contents {
		contentResponses = append(contentResponses, content.ToResponse())
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    contentResponses,
		Total:   total,
		Page:    pagination.Page,
		Limit:   pagination.Limit,
	})
}

// GetContentByID godoc
// @Summary Get content by ID
// @Description Retrieve specific content by its ID
// @Tags content
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {object} models.ContentResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /content/{id} [get]
func GetContentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid content ID",
		})
		return
	}

	var content models.Content
	if err := db.Preload("Category").Preload("Author").
		Where("id = ? AND is_published = ?", id, true).
		First(&content).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Content not found",
		})
		return
	}

	// Increment view count
	db.Model(&content).Update("view_count", content.ViewCount+1)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    content.ToResponse(),
	})
}

// AdminGetContent godoc
// @Summary Get all content (admin)
// @Description Retrieve all content including unpublished (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Router /admin/content [get]
func AdminGetContent(c *gin.Context) {
	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid pagination parameters",
			Error:   err.Error(),
		})
		return
	}

	var contents []models.Content
	var total int64

	query := db.Model(&models.Content{})

	// Count total records
	query.Count(&total)

	// Get paginated results
	if err := query.Preload("Category").Preload("Author").
		Offset(pagination.GetOffset()).
		Limit(pagination.Limit).
		Order("created_at DESC").
		Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch content",
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	var contentResponses []models.ContentResponse
	for _, content := range contents {
		contentResponses = append(contentResponses, content.ToResponse())
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    contentResponses,
		Total:   total,
		Page:    pagination.Page,
		Limit:   pagination.Limit,
	})
}

// CreateContent godoc
// @Summary Create new content
// @Description Create new content (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateContentRequest true "Content details"
// @Success 201 {object} models.ContentResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /admin/content [post]
func CreateContent(c *gin.Context) {
	var req models.CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("userID")

	content := models.Content{
		Title:       req.Title,
		Description: req.Description,
		Body:        req.Body,
		CategoryID:  req.CategoryID,
		AuthorID:    userID.(uint),
		Tags:        req.Tags,
		ImageURL:    req.ImageURL,
		IsPublished: req.IsPublished,
		ViewCount:   0,
	}

	if err := db.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create content",
			Error:   err.Error(),
		})
		return
	}

	// Load relationships for response
	db.Preload("Category").Preload("Author").First(&content, content.ID)

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Data:    content.ToResponse(),
	})
}

// UpdateContent godoc
// @Summary Update content
// @Description Update existing content (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Content ID"
// @Param request body models.UpdateContentRequest true "Content update details"
// @Success 200 {object} models.ContentResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/content/{id} [put]
func UpdateContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid content ID",
		})
		return
	}

	var content models.Content
	if err := db.First(&content, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Content not found",
		})
		return
	}

	var req models.UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Update fields if provided
	if req.Title != "" {
		content.Title = req.Title
	}
	if req.Description != "" {
		content.Description = req.Description
	}
	if req.Body != "" {
		content.Body = req.Body
	}
	if req.CategoryID != 0 {
		content.CategoryID = req.CategoryID
	}
	if req.Tags != "" {
		content.Tags = req.Tags
	}
	if req.ImageURL != "" {
		content.ImageURL = req.ImageURL
	}
	if req.IsPublished != nil {
		content.IsPublished = *req.IsPublished
	}

	if err := db.Save(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update content",
			Error:   err.Error(),
		})
		return
	}

	// Load relationships for response
	db.Preload("Category").Preload("Author").First(&content, content.ID)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    content.ToResponse(),
	})
}

// DeleteContent godoc
// @Summary Delete content
// @Description Delete content (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Content ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/content/{id} [delete]
func DeleteContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid content ID",
		})
		return
	}

	if err := db.Delete(&models.Content{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete content",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Content deleted successfully",
	})
}
