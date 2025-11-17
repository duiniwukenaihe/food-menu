package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"food-ordering/config"
	"food-ordering/database"
	"food-ordering/handlers"
	"food-ordering/middleware"
	"food-ordering/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost/food_ordering_test?sslmode=disable"
	}

	var err error
	testDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}

	err = testDB.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to ping test database: %v", err))
	}

	err = initializeTestDatabase()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize test database: %v", err))
	}

	seedTestData()
}

func teardown() {
	cleanupTestDatabase()
	if testDB != nil {
		testDB.Close()
	}
}

func initializeTestDatabase() error {
	schema := `
		-- Drop existing tables if they exist
		DROP TABLE IF EXISTS order_items CASCADE;
		DROP TABLE IF EXISTS orders CASCADE;
		DROP TABLE IF EXISTS user_favorites CASCADE;
		DROP TABLE IF EXISTS recommendation_dishes CASCADE;
		DROP TABLE IF EXISTS recommendations CASCADE;
		DROP TABLE IF EXISTS dish_nutrition CASCADE;
		DROP TABLE IF EXISTS dishes CASCADE;
		DROP TABLE IF EXISTS categories CASCADE;
		DROP TABLE IF EXISTS system_config CASCADE;
		DROP TABLE IF EXISTS users CASCADE;

		-- Users table
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			email VARCHAR(100),
			role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('user', 'admin')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Categories table
		CREATE TABLE categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Dishes table
		CREATE TABLE dishes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			category_id INTEGER REFERENCES categories(id),
			price DECIMAL(10,2) NOT NULL,
			image_url VARCHAR(500),
			video_url VARCHAR(500),
			cooking_steps TEXT,
			is_seasonal BOOLEAN DEFAULT FALSE,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Dish nutrition table
		CREATE TABLE dish_nutrition (
			id SERIAL PRIMARY KEY,
			dish_id INTEGER REFERENCES dishes(id) ON DELETE CASCADE,
			calories INTEGER,
			protein DECIMAL(5,2),
			fat DECIMAL(5,2),
			carbohydrates DECIMAL(5,2),
			fiber DECIMAL(5,2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Orders table
		CREATE TABLE orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			total_amount DECIMAL(10,2) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'preparing', 'ready', 'completed', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Order items table
		CREATE TABLE order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
			dish_id INTEGER REFERENCES dishes(id),
			quantity INTEGER NOT NULL DEFAULT 1,
			price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Recommendations table
		CREATE TABLE recommendations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			meat_count INTEGER DEFAULT 1,
			vegetable_count INTEGER DEFAULT 2,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Recommendation dishes table
		CREATE TABLE recommendation_dishes (
			id SERIAL PRIMARY KEY,
			recommendation_id INTEGER REFERENCES recommendations(id) ON DELETE CASCADE,
			dish_id INTEGER REFERENCES dishes(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(recommendation_id, dish_id)
		);

		-- User favorites table
		CREATE TABLE user_favorites (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			dish_id INTEGER REFERENCES dishes(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, dish_id)
		);

		-- System config table
		CREATE TABLE system_config (
			id SERIAL PRIMARY KEY,
			config_key VARCHAR(100) UNIQUE NOT NULL,
			config_value TEXT,
			description TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Create indexes
		CREATE INDEX idx_users_username ON users(username);
		CREATE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_dishes_category ON dishes(category_id);
		CREATE INDEX idx_dishes_seasonal ON dishes(is_seasonal);
		CREATE INDEX idx_orders_user ON orders(user_id);
		CREATE INDEX idx_orders_status ON orders(status);
		CREATE INDEX idx_order_items_order ON order_items(order_id);
	`

	_, err := testDB.Exec(schema)
	return err
}

func seedTestData() {
	// Create categories
	_, _ = testDB.Exec(`
		INSERT INTO categories (name, description) VALUES 
		('肉类', '各种肉类菜品'),
		('蔬菜类', '新鲜蔬菜菜品'),
		('汤类', '营养汤品'),
		('主食', '米饭面食'),
		('甜品', '餐后甜点'),
		('饮品', '各种饮料')
	`)

	// Hash password for test users
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	testUserHash, _ := bcrypt.GenerateFromPassword([]byte("testuser123"), bcrypt.DefaultCost)

	// Create test users
	_, _ = testDB.Exec(`
		INSERT INTO users (username, password_hash, email, role, created_at, updated_at) VALUES 
		('admin', $1, 'admin@example.com', 'admin', NOW(), NOW()),
		('user', $2, 'user@example.com', 'user', NOW(), NOW()),
		('testuser', $3, 'testuser@example.com', 'user', NOW(), NOW())
	`, string(adminHash), string(userHash), string(testUserHash))

	// Create test dishes
	_, _ = testDB.Exec(`
		INSERT INTO dishes (name, description, category_id, price, image_url, cooking_steps, is_seasonal, is_active, created_at, updated_at) VALUES
		('红烧肉', '经典红烧肉', 1, 28.00, 'https://example.com/hongshaorou.jpg', '1. 准备食材\n2. 炒糖色\n3. 红烧', FALSE, TRUE, NOW(), NOW()),
		('清炒时蔬', '新鲜时令蔬菜', 2, 18.00, 'https://example.com/shishu.jpg', '1. 清洗蔬菜\n2. 热锅下油\n3. 清炒', FALSE, TRUE, NOW(), NOW()),
		('冬瓜汤', '清凉冬瓜汤', 3, 12.00, 'https://example.com/donggua.jpg', '1. 准备冬瓜\n2. 清汤煮\n3. 出锅', FALSE, TRUE, NOW(), NOW()),
		('米饭', '白米饭', 4, 3.00, 'https://example.com/mifan.jpg', '1. 洗米\n2. 加水\n3. 蒸熟', FALSE, TRUE, NOW(), NOW()),
		('提拉米苏', '经典提拉米苏', 5, 35.00, 'https://example.com/tiramisu.jpg', '1. 浸泡手指饼\n2. 混合芝士\n3. 冷藏', TRUE, TRUE, NOW(), NOW()),
		('豆浆', '热豆浆', 6, 6.00, 'https://example.com/doujiang.jpg', '1. 泡豆\n2. 研磨\n3. 煮沸', FALSE, TRUE, NOW(), NOW())
	`)

	// Create system config
	_, _ = testDB.Exec(`
		INSERT INTO system_config (config_key, config_value, description, updated_at) VALUES
		('default_meat_count', '1', '默认荤菜数量', NOW()),
		('default_vegetable_count', '2', '默认素菜数量', NOW()),
		('max_dish_count', '6', '菜品最大数量', NOW()),
		('s3_endpoint', '', 'S3端点', NOW()),
		('s3_access_key', '', 'S3访问密钥', NOW()),
		('s3_secret_key', '', 'S3密钥', NOW()),
		('s3_bucket', '', 'S3存储桶', NOW()),
		('s3_region', '', 'S3区域', NOW())
	`)

	// Create recommendations
	_, _ = testDB.Exec(`
		INSERT INTO recommendations (name, description, meat_count, vegetable_count, is_active, created_at) VALUES
		('经典搭配', '一荤两素的经典搭配', 1, 2, TRUE, NOW()),
		('丰盛套餐', '两荤两素的丰盛搭配', 2, 2, TRUE, NOW()),
		('素食套餐', '三素一汤的健康搭配', 0, 3, TRUE, NOW())
	`)
}

func cleanupTestDatabase() {
	testDB.Exec(`
		DROP TABLE IF EXISTS order_items CASCADE;
		DROP TABLE IF EXISTS orders CASCADE;
		DROP TABLE IF EXISTS user_favorites CASCADE;
		DROP TABLE IF EXISTS recommendation_dishes CASCADE;
		DROP TABLE IF EXISTS recommendations CASCADE;
		DROP TABLE IF EXISTS dish_nutrition CASCADE;
		DROP TABLE IF EXISTS dishes CASCADE;
		DROP TABLE IF EXISTS categories CASCADE;
		DROP TABLE IF EXISTS system_config CASCADE;
		DROP TABLE IF EXISTS users CASCADE;
	`)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	cfg := &config.Config{
		Port:       "8080",
		JWTSecret:  "your-secret-key",
		DatabaseURL: "postgres://postgres:password@localhost/food_ordering_test?sslmode=disable",
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := handlers.NewHandler(testDB, cfg)

	api := r.Group("/api/v1")
	{
		public := api.Group("/")
		{
			public.POST("/login", handler.Login)
			public.GET("/dishes", handler.GetDishes)
			public.GET("/dishes/:id", handler.GetDish)
			public.GET("/categories", handler.GetCategories)
			public.GET("/recommendations", handler.GetRecommendations)
			public.GET("/seasonal-dishes", handler.GetSeasonalDishes)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", handler.GetProfile)
			protected.POST("/orders", handler.CreateOrder)
			protected.GET("/orders", handler.GetOrders)
			protected.POST("/favorites/:dishId", handler.AddToFavorites)
			protected.DELETE("/favorites/:dishId", handler.RemoveFromFavorites)
			protected.GET("/favorites", handler.GetFavorites)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", handler.GetUsers)
			admin.POST("/dishes", handler.CreateDish)
			admin.PUT("/dishes/:id", handler.UpdateDish)
			admin.DELETE("/dishes/:id", handler.DeleteDish)
			admin.POST("/categories", handler.CreateCategory)
			admin.PUT("/categories/:id", handler.UpdateCategory)
			admin.DELETE("/categories/:id", handler.DeleteCategory)
			admin.GET("/config", handler.GetConfig)
			admin.PUT("/config", handler.UpdateConfig)
		}
	}

	return r
}

func getAuthToken(username, password string) (string, error) {
	req := models.LoginRequest{
		Username: username,
		Password: password,
	}

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	setupRouter().ServeHTTP(w, request)

	var response models.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Token, nil
}

// Test Login - Happy Path
func TestLoginSuccess(t *testing.T) {
	r := setupRouter()

	payload := []byte(`{"username":"admin","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Token == "" {
		t.Error("Token should not be empty")
	}

	if response.User.Username != "admin" {
		t.Errorf("Expected username admin, got %s", response.User.Username)
	}

	if response.User.Role != "admin" {
		t.Errorf("Expected role admin, got %s", response.User.Role)
	}
}

// Test Login - Invalid credentials
func TestLoginInvalidCredentials(t *testing.T) {
	r := setupRouter()

	payload := []byte(`{"username":"admin","password":"wrongpassword"}`)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Test Login - Non-existent user
func TestLoginNonExistentUser(t *testing.T) {
	r := setupRouter()

	payload := []byte(`{"username":"nonexistent","password":"password123"}`)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Test Get Dishes
func TestGetDishes(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/dishes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if dishes, ok := response["dishes"]; !ok || dishes == nil {
		t.Error("Response should contain dishes array")
	}

	if total, ok := response["total"]; !ok || total == nil {
		t.Error("Response should contain total count")
	}
}

// Test Get Single Dish
func TestGetDish(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/dishes/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var dish models.Dish
	err := json.Unmarshal(w.Body.Bytes(), &dish)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if dish.ID != 1 {
		t.Errorf("Expected dish ID 1, got %d", dish.ID)
	}
}

// Test Get Non-existent Dish
func TestGetNonExistentDish(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/dishes/99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// Test Get Categories
func TestGetCategories(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/categories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var categories []models.Category
	err := json.Unmarshal(w.Body.Bytes(), &categories)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(categories) == 0 {
		t.Error("Should return at least one category")
	}
}

// Test Get Recommendations
func TestGetRecommendations(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/recommendations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var recommendations []models.Recommendation
	err := json.Unmarshal(w.Body.Bytes(), &recommendations)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(recommendations) == 0 {
		t.Error("Should return at least one recommendation")
	}
}

// Test Get Seasonal Dishes
func TestGetSeasonalDishes(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/seasonal-dishes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var dishes []models.Dish
	err := json.Unmarshal(w.Body.Bytes(), &dishes)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	for _, dish := range dishes {
		if !dish.IsSeasonal {
			t.Error("All returned dishes should be seasonal")
		}
	}
}

// Test Get Profile - Unauthorized
func TestGetProfileUnauthorized(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Test Get Profile - Authorized
func TestGetProfileAuthorized(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")
	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var user models.User
	json.Unmarshal(w.Body.Bytes(), &user)

	if user.Username != "admin" {
		t.Errorf("Expected username admin, got %s", user.Username)
	}
}

// Test Create Order - Success
func TestCreateOrderSuccess(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	orderReq := models.CreateOrderRequest{
		Items: []models.CreateOrderItemRequest{
			{DishID: 1, Quantity: 2},
			{DishID: 2, Quantity: 1},
		},
	}

	payload, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Response: %s", w.Code, w.Body.String())
	}

	var order models.Order
	json.Unmarshal(w.Body.Bytes(), &order)

	if order.ID == 0 {
		t.Error("Order ID should be set")
	}

	if order.Status != "pending" {
		t.Errorf("Expected status pending, got %s", order.Status)
	}

	if len(order.Items) != 2 {
		t.Errorf("Expected 2 order items, got %d", len(order.Items))
	}
}

// Test Create Order - Invalid (empty items)
func TestCreateOrderInvalidEmpty(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	orderReq := models.CreateOrderRequest{
		Items: []models.CreateOrderItemRequest{},
	}

	payload, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test Create Order - Non-existent Dish
func TestCreateOrderNonExistentDish(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	orderReq := models.CreateOrderRequest{
		Items: []models.CreateOrderItemRequest{
			{DishID: 99999, Quantity: 1},
		},
	}

	payload, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test Get Orders
func TestGetOrders(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")
	req := httptest.NewRequest("GET", "/api/v1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["orders"]; !ok {
		t.Error("Response should contain orders array")
	}
}

// Test Add to Favorites - Success
func TestAddToFavoritesSuccess(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")
	req := httptest.NewRequest("POST", "/api/v1/favorites/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}
}

// Test Add to Favorites - Duplicate
func TestAddToFavoritesDuplicate(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	req1 := httptest.NewRequest("POST", "/api/v1/favorites/1", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First add failed with status %d", w1.Code)
	}

	req2 := httptest.NewRequest("POST", "/api/v1/favorites/1", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate, got %d", w2.Code)
	}
}

// Test Add to Favorites - Non-existent Dish
func TestAddToFavoritesNonExistentDish(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")
	req := httptest.NewRequest("POST", "/api/v1/favorites/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// Test Remove from Favorites - Success
func TestRemoveFromFavoritesSuccess(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	// First add to favorites
	req1 := httptest.NewRequest("POST", "/api/v1/favorites/2", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// Then remove
	req2 := httptest.NewRequest("DELETE", "/api/v1/favorites/2", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}
}

// Test Remove from Favorites - Non-existent
func TestRemoveFromFavoritesNonExistent(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")
	req := httptest.NewRequest("DELETE", "/api/v1/favorites/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// Test Get Favorites
func TestGetFavorites(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	req := httptest.NewRequest("POST", "/api/v1/favorites/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	req2 := httptest.NewRequest("GET", "/api/v1/favorites", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response)

	if _, ok := response["favorites"]; !ok {
		t.Error("Response should contain favorites array")
	}
}

// Test Admin Get Users
func TestAdminGetUsers(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")
	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["users"]; !ok {
		t.Error("Response should contain users array")
	}
}

// Test Admin Create Dish
func TestAdminCreateDish(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")

	dishReq := models.CreateDishRequest{
		Name:         "测试菜品",
		Description:  "这是一个测试菜品",
		CategoryID:   1,
		Price:        25.50,
		ImageURL:     "https://example.com/test.jpg",
		CookingSteps: "1. 准备\n2. 烹饪",
		IsSeasonal:   false,
	}

	payload, _ := json.Marshal(dishReq)
	req := httptest.NewRequest("POST", "/api/v1/admin/dishes", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Response: %s", w.Code, w.Body.String())
	}

	var dish models.Dish
	json.Unmarshal(w.Body.Bytes(), &dish)

	if dish.Name != "测试菜品" {
		t.Errorf("Expected name 测试菜品, got %s", dish.Name)
	}

	if dish.Price != 25.50 {
		t.Errorf("Expected price 25.50, got %f", dish.Price)
	}
}

// Test Admin Create Dish - Unauthorized
func TestAdminCreateDishUnauthorized(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	dishReq := models.CreateDishRequest{
		Name:       "测试菜品",
		CategoryID: 1,
		Price:      25.50,
	}

	payload, _ := json.Marshal(dishReq)
	req := httptest.NewRequest("POST", "/api/v1/admin/dishes", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

// Test Admin Update Dish
func TestAdminUpdateDish(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")

	newName := "更新的菜品名称"
	newPrice := 30.00

	dishReq := models.UpdateDishRequest{
		Name:  &newName,
		Price: &newPrice,
	}

	payload, _ := json.Marshal(dishReq)
	req := httptest.NewRequest("PUT", "/api/v1/admin/dishes/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var dish models.Dish
	json.Unmarshal(w.Body.Bytes(), &dish)

	if dish.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, dish.Name)
	}

	if dish.Price != newPrice {
		t.Errorf("Expected price %f, got %f", newPrice, dish.Price)
	}
}

// Test Admin Delete Dish
func TestAdminDeleteDish(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")
	req := httptest.NewRequest("DELETE", "/api/v1/admin/dishes/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify dish is soft-deleted
	req2 := httptest.NewRequest("GET", "/api/v1/dishes/1", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Deleted dish should not be found, got status %d", w2.Code)
	}
}

// Test Admin Create Category
func TestAdminCreateCategory(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")

	catReq := models.CreateCategoryRequest{
		Name:        "新分类",
		Description: "这是一个新的分类",
	}

	payload, _ := json.Marshal(catReq)
	req := httptest.NewRequest("POST", "/api/v1/admin/categories", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var category models.Category
	json.Unmarshal(w.Body.Bytes(), &category)

	if category.Name != "新分类" {
		t.Errorf("Expected name 新分类, got %s", category.Name)
	}
}

// Test Admin Update Category
func TestAdminUpdateCategory(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")

	newName := "更新的分类"
	catReq := models.UpdateCategoryRequest{
		Name: &newName,
	}

	payload, _ := json.Marshal(catReq)
	req := httptest.NewRequest("PUT", "/api/v1/admin/categories/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var category models.Category
	json.Unmarshal(w.Body.Bytes(), &category)

	if category.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, category.Name)
	}
}

// Test Admin Get Config
func TestAdminGetConfig(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")
	req := httptest.NewRequest("GET", "/api/v1/admin/config", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var configs []models.SystemConfig
	err := json.Unmarshal(w.Body.Bytes(), &configs)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(configs) == 0 {
		t.Error("Should return at least one config")
	}
}

// Test Admin Update Config
func TestAdminUpdateConfig(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")

	configMap := map[string]string{
		"default_meat_count": "2",
	}

	payload, _ := json.Marshal(configMap)
	req := httptest.NewRequest("PUT", "/api/v1/admin/config", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// Test Missing Authorization Header
func TestMissingAuthorizationHeader(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	body, _ := ioutil.ReadAll(w.Body)
	if !bytes.Contains(body, []byte("Authorization header required")) {
		t.Error("Error message should mention missing Authorization header")
	}
}

// Test Invalid Token
func TestInvalidToken(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_here")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Test Admin-only endpoint with regular user
func TestAdminEndpointRegularUser(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")
	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

// Test Create Order - DB Side Effects
func TestCreateOrderDatabaseSideEffects(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("testuser", "testuser123")

	// Get user ID
	var userID int
	testDB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&userID)

	orderReq := models.CreateOrderRequest{
		Items: []models.CreateOrderItemRequest{
			{DishID: 1, Quantity: 2},
		},
	}

	payload, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var orderResponse models.Order
	json.Unmarshal(w.Body.Bytes(), &orderResponse)

	// Verify order exists in database
	var orderCount int
	testDB.QueryRow("SELECT COUNT(*) FROM orders WHERE id = $1 AND user_id = $2", orderResponse.ID, userID).Scan(&orderCount)

	if orderCount != 1 {
		t.Error("Order should exist in database")
	}

	// Verify order items exist in database
	var itemCount int
	testDB.QueryRow("SELECT COUNT(*) FROM order_items WHERE order_id = $1", orderResponse.ID).Scan(&itemCount)

	if itemCount != 1 {
		t.Errorf("Expected 1 order item, got %d", itemCount)
	}
}

// Test Add to Favorites - DB Side Effects
func TestAddToFavoritesDatabaseSideEffects(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("testuser", "testuser123")

	var userID int
	testDB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&userID)

	req := httptest.NewRequest("POST", "/api/v1/favorites/3", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify favorite exists in database
	var favCount int
	testDB.QueryRow("SELECT COUNT(*) FROM user_favorites WHERE user_id = $1 AND dish_id = 3", userID).Scan(&favCount)

	if favCount != 1 {
		t.Error("Favorite should exist in database")
	}
}

// Test Dish Filtering by Category
func TestGetDishesByCategory(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/dishes?category_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if dishesRaw, ok := response["dishes"]; ok {
		dishesBytes, _ := json.Marshal(dishesRaw)
		var dishes []models.Dish
		json.Unmarshal(dishesBytes, &dishes)

		for _, dish := range dishes {
			if dish.CategoryID != 1 {
				t.Errorf("Expected category 1, got %d", dish.CategoryID)
			}
		}
	}
}

// Test Pagination
func TestPagination(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/dishes?page=1&limit=2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if page, ok := response["page"]; ok {
		if int(page.(float64)) != 1 {
			t.Errorf("Expected page 1, got %v", page)
		}
	}

	if limit, ok := response["limit"]; ok {
		if int(limit.(float64)) != 2 {
			t.Errorf("Expected limit 2, got %v", limit)
		}
	}
}

// Test Dish Search
func TestDishSearch(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/dishes?search=红烧肉", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if dishesRaw, ok := response["dishes"]; ok {
		dishesBytes, _ := json.Marshal(dishesRaw)
		var dishes []models.Dish
		json.Unmarshal(dishesBytes, &dishes)

		if len(dishes) > 0 {
			found := false
			for _, dish := range dishes {
				if dish.Name == "红烧肉" {
					found = true
					break
				}
			}
			if !found {
				t.Error("Should find 红烧肉 in search results")
			}
		}
	}
}

// Test Response Time
func TestResponseTime(t *testing.T) {
	r := setupRouter()

	start := time.Now()
	req := httptest.NewRequest("GET", "/api/v1/dishes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	elapsed := time.Since(start)

	if elapsed > 1*time.Second {
		t.Logf("Warning: Response took %v, which is longer than expected", elapsed)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// Test Order With Invalid Quantity
func TestCreateOrderInvalidQuantity(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	orderReq := models.CreateOrderRequest{
		Items: []models.CreateOrderItemRequest{
			{DishID: 1, Quantity: 0},
		},
	}

	payload, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test User Registration Endpoint - if exists
func TestLoginWithEmptyCredentials(t *testing.T) {
	r := setupRouter()

	payload := []byte(`{"username":"","password":""}`)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test Admin Delete Category with Active Dishes
func TestAdminDeleteCategoryWithActiveDishes(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")

	// Get a category that has active dishes
	var categoryID int
	testDB.QueryRow(`
		SELECT c.id FROM categories c
		WHERE EXISTS (SELECT 1 FROM dishes d WHERE d.category_id = c.id AND d.is_active = true)
		LIMIT 1
	`).Scan(&categoryID)

	if categoryID > 0 {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/admin/categories/%d", categoryID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	}
}

// Test Bearer Token Format
func TestInvalidBearerTokenFormat(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Authorization", "InvalidToken")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Test Favorites Pagination
func TestFavoritesPagination(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	req := httptest.NewRequest("GET", "/api/v1/favorites?page=1&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["favorites"]; !ok {
		t.Error("Response should contain favorites array")
	}
}

// Test Orders Pagination
func TestOrdersPagination(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("user", "user123")

	req := httptest.NewRequest("GET", "/api/v1/orders?page=1&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if _, ok := response["orders"]; !ok {
		t.Error("Response should contain orders array")
	}
}

// Test Admin Users Search
func TestAdminUsersSearch(t *testing.T) {
	r := setupRouter()

	token, _ := getAuthToken("admin", "admin123")

	req := httptest.NewRequest("GET", "/api/v1/admin/users?search=admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if usersRaw, ok := response["users"]; ok {
		usersBytes, _ := json.Marshal(usersRaw)
		var users []models.User
		json.Unmarshal(usersBytes, &users)

		for _, user := range users {
			if !bytes.Contains([]byte(user.Username), []byte("admin")) && !bytes.Contains([]byte(user.Email), []byte("admin")) {
				t.Errorf("Search result should match query: %s", user.Username)
			}
		}
	}
}
