//go:build integration
// +build integration

package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "example.com/app/internal/api"
    "example.com/app/internal/database"
    "example.com/app/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    db, err := database.InitDB()
    if err != nil {
        panic("Failed to connect to test database")
    }
    
    // Clean database
    db.Exec("DELETE FROM recommendations")
    db.Exec("DELETE FROM content")
    db.Exec("DELETE FROM categories")
    db.Exec("DELETE FROM users")
    
    return db
}

func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    
    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })
    
    v1 := r.Group("/api/v1")
    {
        v1.POST("/auth/register", api.Register)
        v1.POST("/auth/login", api.Login)
        v1.GET("/content", api.GetContent)
        v1.GET("/content/:id", api.GetContentByID)
        v1.GET("/categories", api.GetCategories)
        
        protected := v1.Group("/")
        protected.Use(api.AuthMiddleware())
        {
            protected.GET("/auth/profile", api.GetProfile)
            protected.GET("/recommendations", api.GetRecommendations)
        }
    }
    
    return r
}

func TestUserRegistrationAndLogin(t *testing.T) {
    db := setupTestDB()
    api.InitDatabase(db)
    router := setupRouter()
    
    // Test user registration
    registerReq := models.RegisterRequest{
        Username:  "testuser",
        Email:     "test@example.com",
        Password:  "password123",
        FirstName: "Test",
        LastName:  "User",
    }
    
    registerBody, _ := json.Marshal(registerReq)
    req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var authResponse models.AuthResponse
    err := json.Unmarshal(w.Body.Bytes(), &authResponse)
    assert.NoError(t, err)
    assert.NotEmpty(t, authResponse.Token)
    assert.Equal(t, registerReq.Username, authResponse.User.Username)
    
    // Test user login
    loginReq := models.LoginRequest{
        Username: "testuser",
        Password: "password123",
    }
    
    loginBody, _ := json.Marshal(loginReq)
    req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
    req.Header.Set("Content-Type", "application/json")
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    err = json.Unmarshal(w.Body.Bytes(), &authResponse)
    assert.NoError(t, err)
    assert.NotEmpty(t, authResponse.Token)
    assert.Equal(t, loginReq.Username, authResponse.User.Username)
}

func TestContentRetrieval(t *testing.T) {
    db := setupTestDB()
    api.InitDatabase(db)
    router := setupRouter()
    
    // Create test data
    category := models.Category{
        Name:     "Test Category",
        IsActive: true,
    }
    db.Create(&category)
    
    user := models.User{
        Username:  "testuser",
        Email:     "test@example.com",
        Password:  "hashedpassword",
        FirstName: "Test",
        LastName:  "User",
        Role:      "user",
        IsActive:  true,
    }
    db.Create(&user)
    
    content := models.Content{
        Title:       "Test Article",
        Description: "Test description",
        Body:        "Test content body",
        CategoryID:  category.ID,
        AuthorID:    user.ID,
        IsPublished: true,
        ViewCount:   0,
    }
    db.Create(&content)
    
    // Test get all content
    req, _ := http.NewRequest("GET", "/api/v1/content", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var paginatedResponse models.PaginatedResponse
    err := json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
    assert.NoError(t, err)
    assert.Equal(t, int64(1), paginatedResponse.Total)
    
    // Test get content by ID
    req, _ = http.NewRequest("GET", "/api/v1/content/1", nil)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var successResponse models.SuccessResponse
    err = json.Unmarshal(w.Body.Bytes(), &successResponse)
    assert.NoError(t, err)
    
    contentData, ok := successResponse.Data.(map[string]interface{})
    assert.True(t, ok)
    assert.Equal(t, "Test Article", contentData["title"])
}

func TestProtectedRoutes(t *testing.T) {
    db := setupTestDB()
    api.InitDatabase(db)
    router := setupRouter()
    
    // Test accessing protected route without token
    req, _ := http.NewRequest("GET", "/api/v1/auth/profile", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    
    // Create and login user to get token
    registerReq := models.RegisterRequest{
        Username:  "testuser",
        Email:     "test@example.com",
        Password:  "password123",
        FirstName: "Test",
        LastName:  "User",
    }
    
    registerBody, _ := json.Marshal(registerReq)
    req, _ = http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
    req.Header.Set("Content-Type", "application/json")
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    var authResponse models.AuthResponse
    json.Unmarshal(w.Body.Bytes(), &authResponse)
    
    // Test accessing protected route with valid token
    req, _ = http.NewRequest("GET", "/api/v1/auth/profile", nil)
    req.Header.Set("Authorization", "Bearer "+authResponse.Token)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
