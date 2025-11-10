package api

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    "example.com/app/internal/models"
)

// GetDishes godoc
// @Summary Get all dishes (public)
// @Description Get a paginated list of dishes with optional filters including category
// @Tags dishes
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name or description"
// @Param isSeasonal query bool false "Filter by seasonal dishes"
// @Param isActive query bool false "Filter by active status"
// @Param tags query string false "Filter by tags (comma-separated)"
// @Param categoryId query int false "Filter by category ID"
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /dishes [get]
func GetDishes(c *gin.Context) {
    var params models.DishQueryParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid query parameters",
            Error:   err.Error(),
        })
        return
    }

    var dishes []models.Dish
    var total int64

    query := db.Model(&models.Dish{})

    // Apply filters
    if params.Search != "" {
        query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
    }

    if params.IsSeasonal != nil {
        query = query.Where("is_seasonal = ?", *params.IsSeasonal)
    }

    if params.IsActive != nil {
        query = query.Where("is_active = ?", *params.IsActive)
    } else {
        // By default, only show active dishes for public endpoint
        query = query.Where("is_active = ?", true)
    }

    if params.CategoryID > 0 {
        query = query.Where("category_id = ?", params.CategoryID)
    }

    if params.Tags != "" {
        tags := strings.Split(params.Tags, ",")
        for _, tag := range tags {
            query = query.Where("tags ILIKE ?", "%"+strings.TrimSpace(tag)+"%")
        }
    }

    // Get total count
    query.Count(&total)

    // Apply pagination
    if params.Limit <= 0 {
        params.Limit = 10
    }
    if params.Page <= 0 {
        params.Page = 1
    }

    offset := params.GetOffset()
    if err := query.Preload("Category").Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&dishes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to fetch dishes",
            Error:   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, models.PaginatedResponse{
        Success: true,
        Data:    dishes,
        Total:   total,
        Page:    params.Page,
        Limit:   params.Limit,
    })
}

// GetDishByID godoc
// @Summary Get dish by ID
// @Description Get a single dish by its ID with all details including steps and media
// @Tags dishes
// @Accept json
// @Produce json
// @Param id path int true "Dish ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /dishes/{id} [get]
func GetDishByID(c *gin.Context) {
    id := c.Param("id")

    var dish models.Dish
    if err := db.Preload("Category").Where("id = ? AND is_active = ?", id, true).First(&dish).Error; err != nil {
        c.JSON(http.StatusNotFound, models.ErrorResponse{
            Success: false,
            Message: "Dish not found",
        })
        return
    }

    c.JSON(http.StatusOK, models.SuccessResponse{
        Success: true,
        Data:    dish,
    })
}

// AdminGetDishes godoc
// @Summary Get all dishes (admin)
// @Description Get a paginated list of all dishes including inactive ones
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name or description"
// @Param isSeasonal query bool false "Filter by seasonal dishes"
// @Param isActive query bool false "Filter by active status"
// @Param tags query string false "Filter by tags (comma-separated)"
// @Param categoryId query int false "Filter by category ID"
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/dishes [get]
func AdminGetDishes(c *gin.Context) {
    var params models.DishQueryParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid query parameters",
            Error:   err.Error(),
        })
        return
    }

    var dishes []models.Dish
    var total int64

    query := db.Model(&models.Dish{})

    // Apply filters
    if params.Search != "" {
        query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
    }

    if params.IsSeasonal != nil {
        query = query.Where("is_seasonal = ?", *params.IsSeasonal)
    }

    if params.IsActive != nil {
        query = query.Where("is_active = ?", *params.IsActive)
    }

    if params.CategoryID > 0 {
        query = query.Where("category_id = ?", params.CategoryID)
    }

    if params.Tags != "" {
        tags := strings.Split(params.Tags, ",")
        for _, tag := range tags {
            query = query.Where("tags ILIKE ?", "%"+strings.TrimSpace(tag)+"%")
        }
    }

    // Get total count
    query.Count(&total)

    // Apply pagination
    if params.Limit <= 0 {
        params.Limit = 10
    }
    if params.Page <= 0 {
        params.Page = 1
    }

    offset := params.GetOffset()
    if err := query.Preload("Category").Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&dishes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to fetch dishes",
            Error:   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, models.PaginatedResponse{
        Success: true,
        Data:    dishes,
        Total:   total,
        Page:    params.Page,
        Limit:   params.Limit,
    })
}

// CreateDish godoc
// @Summary Create a new dish
// @Description Create a new dish (admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param dish body models.CreateDishRequest true "Dish data"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/dishes [post]
func CreateDish(c *gin.Context) {
    var req models.CreateDishRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid request body",
            Error:   err.Error(),
        })
        return
    }

    dish := models.Dish{
        Name:            req.Name,
        Description:     req.Description,
        Tags:            req.Tags,
        IsActive:        req.IsActive,
        CategoryID:      req.CategoryID,
        TextSteps:       req.TextSteps,
        IsSeasonal:      req.IsSeasonal,
        AvailableMonths: req.AvailableMonths,
        SeasonalNote:    req.SeasonalNote,
        ImageURL:        req.ImageURL,
        ThumbnailURL:    req.ThumbnailURL,
        GalleryURLs:     req.GalleryURLs,
    }

    if err := db.Create(&dish).Error; err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to create dish",
            Error:   err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, models.SuccessResponse{
        Success: true,
        Message: "Dish created successfully",
        Data:    dish,
    })
}

// UpdateDish godoc
// @Summary Update a dish
// @Description Update an existing dish (admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Dish ID"
// @Param dish body models.UpdateDishRequest true "Dish data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/dishes/{id} [put]
func UpdateDish(c *gin.Context) {
    id := c.Param("id")
    dishID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid dish ID",
        })
        return
    }

    var req models.UpdateDishRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid request body",
            Error:   err.Error(),
        })
        return
    }

    var dish models.Dish
    if err := db.First(&dish, dishID).Error; err != nil {
        c.JSON(http.StatusNotFound, models.ErrorResponse{
            Success: false,
            Message: "Dish not found",
        })
        return
    }

    // Track old media URLs for cleanup
    oldImageURL := dish.ImageURL
    oldThumbnailURL := dish.ThumbnailURL
    oldGalleryURLs := dish.GalleryURLs

    // Update only provided fields
    updates := make(map[string]interface{})
    if req.Name != nil {
        updates["name"] = *req.Name
    }
    if req.Description != nil {
        updates["description"] = *req.Description
    }
    if req.Tags != nil {
        updates["tags"] = *req.Tags
    }
    if req.IsActive != nil {
        updates["is_active"] = *req.IsActive
    }
    if req.CategoryID != nil {
        updates["category_id"] = *req.CategoryID
    }
    if req.TextSteps != nil {
        updates["text_steps"] = *req.TextSteps
    }
    if req.IsSeasonal != nil {
        updates["is_seasonal"] = *req.IsSeasonal
    }
    if req.AvailableMonths != nil {
        updates["available_months"] = *req.AvailableMonths
    }
    if req.SeasonalNote != nil {
        updates["seasonal_note"] = *req.SeasonalNote
    }
    if req.ImageURL != nil {
        updates["image_url"] = *req.ImageURL
    }
    if req.ThumbnailURL != nil {
        updates["thumbnail_url"] = *req.ThumbnailURL
    }
    if req.GalleryURLs != nil {
        updates["gallery_urls"] = *req.GalleryURLs
    }

    if err := db.Model(&dish).Updates(updates).Error; err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to update dish",
            Error:   err.Error(),
        })
        return
    }

    // Cleanup replaced media
    if req.ImageURL != nil && oldImageURL != "" && oldImageURL != *req.ImageURL {
        CleanupReplacedMedia(oldImageURL)
    }
    if req.ThumbnailURL != nil && oldThumbnailURL != "" && oldThumbnailURL != *req.ThumbnailURL {
        CleanupReplacedMedia(oldThumbnailURL)
    }
    if req.GalleryURLs != nil && oldGalleryURLs != "" && oldGalleryURLs != *req.GalleryURLs {
        // Gallery URLs is a JSON array, so we'd need to parse and compare, for now just skip
        // In production, you'd parse the JSON and clean up removed URLs
    }

    // Fetch updated dish
    db.First(&dish, dishID)

    c.JSON(http.StatusOK, models.SuccessResponse{
        Success: true,
        Message: "Dish updated successfully",
        Data:    dish,
    })
}

// DeleteDish godoc
// @Summary Delete a dish
// @Description Delete a dish by ID (admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Dish ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/dishes/{id} [delete]
func DeleteDish(c *gin.Context) {
    id := c.Param("id")

    var dish models.Dish
    if err := db.First(&dish, id).Error; err != nil {
        c.JSON(http.StatusNotFound, models.ErrorResponse{
            Success: false,
            Message: "Dish not found",
        })
        return
    }

    // Cleanup associated media before deleting dish
    if dish.ImageURL != "" {
        CleanupReplacedMedia(dish.ImageURL)
    }
    if dish.ThumbnailURL != "" {
        CleanupReplacedMedia(dish.ThumbnailURL)
    }

    if err := db.Delete(&dish).Error; err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to delete dish",
            Error:   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, models.SuccessResponse{
        Success: true,
        Message: "Dish deleted successfully",
    })
}
