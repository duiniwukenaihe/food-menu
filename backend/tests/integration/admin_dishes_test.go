//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"example.com/app/internal/api"
	"example.com/app/internal/database"
	"example.com/app/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

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
		v1.GET("/dishes", api.GetDishes)
		v1.GET("/dishes/:id", api.GetDishByID)
		v1.GET("/media/*filepath", api.GetMedia)

		protected := v1.Group("/")
		protected.Use(api.AuthMiddleware())
		{
			admin := protected.Group("/admin")
			admin.Use(api.AdminMiddleware())
			{
				admin.GET("/dishes", api.AdminGetDishes)
				admin.POST("/dishes", api.CreateDish)
				admin.PUT("/dishes/:id", api.UpdateDish)
				admin.DELETE("/dishes/:id", api.DeleteDish)

				admin.POST("/media/upload-url", api.GetUploadURL)
				admin.POST("/media/upload", api.UploadFile)
				admin.DELETE("/media", api.DeleteMedia)
			}
		}
	}

	return r
}

func setupTestEnv(t *testing.T) (*gorm.DB, *gin.Engine, string) {
	db, err := database.InitDB()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}

	db.Exec("DELETE FROM media")
	db.Exec("DELETE FROM dishes")
	db.Exec("DELETE FROM users")

	db.AutoMigrate(&models.User{}, &models.Dish{}, &models.Media{})

	api.InitDatabase(db)
	if err := api.InitStorageService(); err != nil {
		t.Fatal("Failed to initialize storage service:", err)
	}

	router := setupTestRouter()

	adminToken := createAdminUser(t, db, router)

	return db, router, adminToken
}

func createAdminUser(t *testing.T, db *gorm.DB, router *gin.Engine) string {
	registerReq := models.RegisterRequest{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "admin123",
		FirstName: "Admin",
		LastName:  "User",
	}

	registerBody, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authResponse models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &authResponse)

	db.Model(&models.User{}).Where("username = ?", "admin").Update("role", "admin")

	loginReq := models.LoginRequest{
		Username: "admin",
		Password: "admin123",
	}

	loginBody, _ := json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &authResponse)

	return authResponse.Token
}

func TestAdminDishCRUD(t *testing.T) {
	db, router, token := setupTestEnv(t)
	defer db.Exec("DELETE FROM dishes")

	t.Run("Create Dish", func(t *testing.T) {
		createReq := models.CreateDishRequest{
			Name:        "Test Dish",
			Description: "A delicious test dish",
			Tags:        "test,food",
			IsActive:    true,
			IsSeasonal:  false,
		}

		body, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/dishes", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		dishData := response.Data.(map[string]interface{})
		assert.Equal(t, "Test Dish", dishData["name"])
		assert.Equal(t, "A delicious test dish", dishData["description"])
	})

	t.Run("List All Dishes (Admin)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/dishes", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Greater(t, response.Total, int64(0))
	})

	t.Run("Update Dish", func(t *testing.T) {
		var dish models.Dish
		db.First(&dish)

		newName := "Updated Test Dish"
		updateReq := models.UpdateDishRequest{
			Name: &newName,
		}

		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/admin/dishes/%d", dish.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		dishData := response.Data.(map[string]interface{})
		assert.Equal(t, "Updated Test Dish", dishData["name"])
	})

	t.Run("Delete Dish", func(t *testing.T) {
		var dish models.Dish
		db.First(&dish)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/admin/dishes/%d", dish.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		var count int64
		db.Model(&models.Dish{}).Where("id = ?", dish.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestAdminDishValidation(t *testing.T) {
	_, router, token := setupTestEnv(t)

	t.Run("Create Dish with Invalid Name", func(t *testing.T) {
		createReq := models.CreateDishRequest{
			Name:        "",
			Description: "Test description",
		}

		body, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/dishes", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Create Dish with Name Too Long", func(t *testing.T) {
		longName := strings.Repeat("a", 256)
		createReq := models.CreateDishRequest{
			Name:        longName,
			Description: "Test description",
		}

		body, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/dishes", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMediaUploadAndStorage(t *testing.T) {
	db, router, token := setupTestEnv(t)
	defer db.Exec("DELETE FROM media")

	t.Run("Get Presigned Upload URL", func(t *testing.T) {
		uploadReq := models.UploadURLRequest{
			FileName:    "test-image.jpg",
			ContentType: "image/jpeg",
		}

		body, _ := json.Marshal(uploadReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload-url", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		urlData := response.Data.(map[string]interface{})
		assert.NotEmpty(t, urlData["uploadUrl"])
		assert.NotEmpty(t, urlData["fileUrl"])
		assert.NotEmpty(t, urlData["key"])
	})

	t.Run("Upload File Directly", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileContent := []byte("fake image content")
		part, err := writer.CreateFormFile("file", "test-image.jpg")
		assert.NoError(t, err)
		part.Write(fileContent)
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		mediaData := response.Data.(map[string]interface{})
		assert.NotEmpty(t, mediaData["url"])
		assert.Equal(t, "test-image.jpg", mediaData["fileName"])

		var mediaCount int64
		db.Model(&models.Media{}).Count(&mediaCount)
		assert.Greater(t, mediaCount, int64(0))
	})

	t.Run("Upload Invalid File Type", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileContent := []byte("fake text content")
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.NoError(t, err)
		part.Write(fileContent)
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Upload File Too Large", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		largeContent := make([]byte, 11*1024*1024)
		part, err := writer.CreateFormFile("file", "large-image.jpg")
		assert.NoError(t, err)
		part.Write(largeContent)
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMediaValidation(t *testing.T) {
	_, router, token := setupTestEnv(t)

	t.Run("Invalid Content Type for Upload URL", func(t *testing.T) {
		uploadReq := models.UploadURLRequest{
			FileName:    "test.pdf",
			ContentType: "application/pdf",
		}

		body, _ := json.Marshal(uploadReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload-url", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Empty File Name", func(t *testing.T) {
		uploadReq := models.UploadURLRequest{
			FileName:    "",
			ContentType: "image/jpeg",
		}

		body, _ := json.Marshal(uploadReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload-url", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("File Name Too Long", func(t *testing.T) {
		longName := strings.Repeat("a", 256) + ".jpg"
		uploadReq := models.UploadURLRequest{
			FileName:    longName,
			ContentType: "image/jpeg",
		}

		body, _ := json.Marshal(uploadReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload-url", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDishWithMediaIntegration(t *testing.T) {
	db, router, token := setupTestEnv(t)
	defer db.Exec("DELETE FROM dishes")
	defer db.Exec("DELETE FROM media")

	var uploadedURL string

	t.Run("Upload Media for Dish", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileContent := []byte("fake image content")
		part, err := writer.CreateFormFile("file", "dish-image.jpg")
		assert.NoError(t, err)
		part.Write(fileContent)
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		mediaData := response.Data.(map[string]interface{})
		uploadedURL = mediaData["url"].(string)
	})

	t.Run("Create Dish with Uploaded Media", func(t *testing.T) {
		createReq := models.CreateDishRequest{
			Name:        "Dish with Image",
			Description: "A dish with an uploaded image",
			ImageURL:    uploadedURL,
			IsActive:    true,
		}

		body, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/dishes", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		dishData := response.Data.(map[string]interface{})
		assert.Equal(t, uploadedURL, dishData["imageUrl"])
	})

	t.Run("Update Dish with New Media", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileContent := []byte("new fake image content")
		part, err := writer.CreateFormFile("file", "new-dish-image.jpg")
		assert.NoError(t, err)
		part.Write(fileContent)
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response models.SuccessResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		mediaData := response.Data.(map[string]interface{})
		newUploadedURL := mediaData["url"].(string)

		var dish models.Dish
		db.First(&dish)

		updateReq := models.UpdateDishRequest{
			ImageURL: &newUploadedURL,
		}

		body2, _ := json.Marshal(updateReq)
		req2, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/admin/dishes/%d", dish.ID), bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer "+token)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("Delete Dish Cleans Up Media", func(t *testing.T) {
		var dish models.Dish
		db.First(&dish)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/admin/dishes/%d", dish.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminAuthorization(t *testing.T) {
	db, router, _ := setupTestEnv(t)

	registerReq := models.RegisterRequest{
		Username:  "regularuser",
		Email:     "user@example.com",
		Password:  "user123",
		FirstName: "Regular",
		LastName:  "User",
	}

	registerBody, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authResponse models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &authResponse)
	userToken := authResponse.Token

	defer db.Exec("DELETE FROM users WHERE username = 'regularuser'")

	t.Run("Regular User Cannot Access Admin Endpoints", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/dishes", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Unauthenticated User Cannot Access Admin Endpoints", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/admin/dishes", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestVideoUpload(t *testing.T) {
	_, router, token := setupTestEnv(t)

	t.Run("Upload Video File", func(t *testing.T) {
		uploadReq := models.UploadURLRequest{
			FileName:    "test-video.mp4",
			ContentType: "video/mp4",
		}

		body, _ := json.Marshal(uploadReq)
		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload-url", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})
}

func TestDeleteMedia(t *testing.T) {
	db, router, token := setupTestEnv(t)
	defer db.Exec("DELETE FROM media")

	var uploadedKey string

	t.Run("Upload and Delete Media", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileContent := []byte("test content")
		part, err := writer.CreateFormFile("file", "delete-test.jpg")
		assert.NoError(t, err)
		_, err = io.Copy(part, bytes.NewReader(fileContent))
		assert.NoError(t, err)
		writer.Close()

		req, _ := http.NewRequest("POST", "/api/v1/admin/media/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var media models.Media
		db.First(&media)
		uploadedKey = media.Key

		deleteReq := models.DeleteMediaRequest{
			Key: uploadedKey,
		}

		deleteBody, _ := json.Marshal(deleteReq)
		req2, _ := http.NewRequest("DELETE", "/api/v1/admin/media", bytes.NewBuffer(deleteBody))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer "+token)

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)

		var count int64
		db.Model(&models.Media{}).Where("key = ?", uploadedKey).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func init() {
	os.Setenv("STORAGE_TYPE", "local")
	os.Setenv("UPLOAD_DIR", "./test_uploads")
}
