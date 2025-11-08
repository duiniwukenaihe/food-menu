package api

import (
	"net/http"
	"strings"
	"example.com/app/internal/auth"
	"example.com/app/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase(database *gorm.DB) {
	db = database
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
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

	// Generate token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Find user
	var user models.User
	if err := db.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/profile [get]
func GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    user.(*models.User).ToResponse(),
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Profile update details"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/profile [put]
func UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
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

	currentUser := user.(*models.User)
	
	// Update fields if provided
	if req.FirstName != "" {
		currentUser.FirstName = req.FirstName
	}
	if req.LastName != "" {
		currentUser.LastName = req.LastName
	}
	if req.Avatar != "" {
		currentUser.Avatar = req.Avatar
	}

	if err := db.Save(currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update profile",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Data:    currentUser.ToResponse(),
	})
}

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid token",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "User not found",
			})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)
		c.Next()
	}
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "User not authenticated",
			})
			c.Abort()
			return
		}

		if userRole != "admin" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Success: false,
				Message: "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
