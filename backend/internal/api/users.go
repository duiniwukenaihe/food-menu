package api

import (
	"net/http"
	"strconv"
	"example.com/app/internal/auth"
	"example.com/app/internal/models"
	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve paginated list of users (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Router /admin/users [get]
func GetUsers(c *gin.Context) {
	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid pagination parameters",
			Error:   err.Error(),
		})
		return
	}

	var users []models.User
	var total int64

	query := db.Model(&models.User{})

	// Count total records
	query.Count(&total)

	// Get paginated results
	if err := query.Offset(pagination.GetOffset()).
		Limit(pagination.Limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    userResponses,
		Total:   total,
		Page:    pagination.Page,
		Limit:   pagination.Limit,
	})
}

// CreateUser godoc
// @Summary Create new user
// @Description Create new user (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RegisterRequest true "User details"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /admin/users [post]
func CreateUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Success: false,
			Message: "User already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to hash password",
		})
		return
	}

	// Create user
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		IsActive:  true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Data:    user.ToResponse(),
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update existing user (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body models.UpdateProfileRequest true "User update details"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    user.ToResponse(),
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
