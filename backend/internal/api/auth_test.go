package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"example.com/app/internal/auth"
	"example.com/app/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of gorm.DB for testing
type MockDB struct {
	mock.Mock
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})
	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Initialize the database for the api package
	InitDatabase(db)
	
	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", Register)
	v1.POST("/auth/login", Login)
	v1.GET("/auth/profile", AuthMiddleware(), GetProfile)
	v1.POST("/auth/logout", AuthMiddleware(), Logout)
	
	return router
}

func TestRegister(t *testing.T) {
	db := setupTestDB()
	router := setupTestRouter(db)

	tests := []struct {
		name           string
		requestBody    models.RegisterRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Successful registration",
			requestBody: models.RegisterRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "Invalid email",
			requestBody: models.RegisterRequest{
				Username:  "testuser2",
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Short password",
			requestBody: models.RegisterRequest{
				Username:  "testuser3",
				Email:     "test3@example.com",
				Password:  "123",
				FirstName: "Test",
				LastName:  "User",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing required fields",
			requestBody: models.RegisterRequest{
				Username:  "testuser4",
				Email:     "test4@example.com",
				Password:  "password123",
				// Missing FirstName and LastName
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectError {
				assert.False(t, response["success"].(bool))
			} else {
				assert.True(t, response["success"].(bool))
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
			}
		})
	}
}

func TestLogin(t *testing.T) {
	db := setupTestDB()
	router := setupTestRouter(db)

	// Create a test user first
	hashedPassword, _ := auth.HashPassword("password123")
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
	}
	db.Create(&user)

	tests := []struct {
		name           string
		requestBody    models.LoginRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Successful login",
			requestBody: models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Invalid username",
			requestBody: models.LoginRequest{
				Username: "wronguser",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "Invalid password",
			requestBody: models.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "Missing fields",
			requestBody: models.LoginRequest{
				Username: "testuser",
				// Missing password
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectError {
				assert.False(t, response["success"].(bool))
			} else {
				assert.True(t, response["success"].(bool))
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "user")
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	db := setupTestDB()
	router := setupTestRouter(db)

	// Create a test user
	hashedPassword, _ := auth.HashPassword("password123")
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
	}
	db.Create(&user)

	// Generate token for the user
	token, _ := auth.GenerateToken(&user)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Successful profile retrieval",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "Invalid authorization format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/auth/profile", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectError {
				assert.False(t, response["success"].(bool))
			} else {
				assert.True(t, response["success"].(bool))
				assert.Contains(t, response, "data")
			}
		})
	}
}

func TestLogout(t *testing.T) {
	db := setupTestDB()
	router := setupTestRouter(db)

	// Create a test user
	hashedPassword, _ := auth.HashPassword("password123")
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
	}
	db.Create(&user)

	// Generate token for the user
	token, _ := auth.GenerateToken(&user)

	t.Run("Successful logout", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Successfully logged out", response["message"])
	})

	t.Run("Logout without token", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
	})
}

func TestAuthMiddleware(t *testing.T) {
	db := setupTestDB()
	
	// Create a test user
	hashedPassword, _ := auth.HashPassword("password123")
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
	}
	db.Create(&user)

	// Generate token for the user
	token, _ := auth.GenerateToken(&user)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		user, exists := c.Get("user")
		if exists {
			c.JSON(http.StatusOK, gin.H{"user": user.(*models.User).Username})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": nil})
		}
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUser   string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
			expectedUser:   "testuser",
		},
		{
			name:           "Missing token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedUser:   "",
		},
		{
			name:           "Invalid token format",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
			expectedUser:   "",
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedUser:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, response["user"])
			}
		})
	}
}