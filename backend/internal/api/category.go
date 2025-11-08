package api

import (
	"net/http"
	"strconv"
	"example.com/app/internal/models"
	"github.com/gin-gonic/gin"
)

// GetCategories godoc
// @Summary Get all active categories
// @Description Retrieve list of active categories
// @Tags categories
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch categories",
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	var categoryResponses []models.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, category.ToResponse())
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    categoryResponses,
	})
}

// AdminGetCategories godoc
// @Summary Get all categories (admin)
// @Description Retrieve all categories including inactive (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse
// @Router /admin/categories [get]
func AdminGetCategories(c *gin.Context) {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch categories",
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	var categoryResponses []models.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, category.ToResponse())
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    categoryResponses,
	})
}

// CreateCategory godoc
// @Summary Create new category
// @Description Create new category (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateCategoryRequest true "Category details"
// @Success 201 {object} models.CategoryResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /admin/categories [post]
func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Check if category already exists
	var existingCategory models.Category
	if err := db.Where("name = ?", req.Name).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Success: false,
			Message: "Category already exists",
		})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    true,
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Data:    category.ToResponse(),
	})
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update existing category (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body models.UpdateCategoryRequest true "Category update details"
// @Success 200 {object} models.CategoryResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid category ID",
		})
		return
	}

	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Category not found",
		})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    category.ToResponse(),
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete category (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid category ID",
		})
		return
	}

	if err := db.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}
