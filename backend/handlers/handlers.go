package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"food-ordering/config"
	"food-ordering/middleware"
	"food-ordering/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db  *sql.DB
	cfg *config.Config
}

func NewHandler(db *sql.DB, cfg *config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// 登录处理
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user models.User
	err := h.db.QueryRow(
		"SELECT id, username, email, role, created_at, updated_at FROM users WHERE username = $1",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 验证密码
	var passwordHash string
	err = h.db.QueryRow("SELECT password_hash FROM users WHERE username = $1", req.Username).Scan(&passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// 获取用户信息
func (h *Handler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var user models.User
	err := h.db.QueryRow(
		"SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// 获取菜品列表
func (h *Handler) GetDishes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	search := c.Query("search")

	offset := (page - 1) * limit

	query := `
		SELECT d.id, d.name, d.description, d.category_id, c.name as category_name,
			   d.price, d.image_url, d.video_url, d.cooking_steps, d.is_seasonal, d.is_active,
			   d.created_at, d.updated_at
		FROM dishes d
		LEFT JOIN categories c ON d.category_id = c.id
		WHERE d.is_active = true
	`
	args := []interface{}{}
	argIndex := 1

	if categoryID > 0 {
		query += fmt.Sprintf(" AND d.category_id = $%d", argIndex)
		args = append(args, categoryID)
		argIndex++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (d.name ILIKE $%d OR d.description ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	query += fmt.Sprintf(" ORDER BY d.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dishes"})
		return
	}
	defer rows.Close()

	var dishes []models.Dish
	for rows.Next() {
		var dish models.Dish
		var categoryName sql.NullString
		
		err := rows.Scan(
			&dish.ID, &dish.Name, &dish.Description, &dish.CategoryID, &categoryName,
			&dish.Price, &dish.ImageURL, &dish.VideoURL, &dish.CookingSteps,
			&dish.IsSeasonal, &dish.IsActive, &dish.CreatedAt, &dish.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan dish"})
			return
		}

		if categoryName.Valid {
			dish.Category = &models.Category{
				ID:   dish.CategoryID,
				Name: categoryName.String,
			}
		}

		dishes = append(dishes, dish)
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM dishes WHERE is_active = true"
	countArgs := []interface{}{}
	countArgIndex := 1

	if categoryID > 0 {
		countQuery += fmt.Sprintf(" AND category_id = $%d", countArgIndex)
		countArgs = append(countArgs, categoryID)
		countArgIndex++
	}

	if search != "" {
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", countArgIndex, countArgIndex+1)
		countArgs = append(countArgs, "%"+search+"%", "%"+search+"%")
	}

	var total int
	h.db.QueryRow(countQuery, countArgs...).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"dishes": dishes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// 获取单个菜品
func (h *Handler) GetDish(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dish ID"})
		return
	}

	var dish models.Dish
	var categoryName sql.NullString

	err = h.db.QueryRow(`
		SELECT d.id, d.name, d.description, d.category_id, c.name as category_name,
			   d.price, d.image_url, d.video_url, d.cooking_steps, d.is_seasonal, d.is_active,
			   d.created_at, d.updated_at
		FROM dishes d
		LEFT JOIN categories c ON d.category_id = c.id
		WHERE d.id = $1
	`, id).Scan(
		&dish.ID, &dish.Name, &dish.Description, &dish.CategoryID, &categoryName,
		&dish.Price, &dish.ImageURL, &dish.VideoURL, &dish.CookingSteps,
		&dish.IsSeasonal, &dish.IsActive, &dish.CreatedAt, &dish.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dish"})
		return
	}

	if categoryName.Valid {
		dish.Category = &models.Category{
			ID:   dish.CategoryID,
			Name: categoryName.String,
		}
	}

	// 获取营养信息
	var nutrition models.DishNutrition
	err = h.db.QueryRow(`
		SELECT id, calories, protein, fat, carbohydrates, fiber, created_at
		FROM dish_nutrition WHERE dish_id = $1
	`, id).Scan(&nutrition.ID, &nutrition.Calories, &nutrition.Protein, &nutrition.Fat, 
		&nutrition.Carbohydrates, &nutrition.Fiber, &nutrition.CreatedAt)

	if err == nil {
		// 可以将营养信息添加到响应中
	}

	c.JSON(http.StatusOK, dish)
}

// 获取分类列表
func (h *Handler) GetCategories(c *gin.Context) {
	rows, err := h.db.Query("SELECT id, name, description, created_at FROM categories ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan category"})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

// 获取推荐列表
func (h *Handler) GetRecommendations(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT r.id, r.name, r.description, r.meat_count, r.vegetable_count, r.is_active, r.created_at
		FROM recommendations r
		WHERE r.is_active = true
		ORDER BY r.created_at
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}
	defer rows.Close()

	var recommendations []models.Recommendation
	for rows.Next() {
		var rec models.Recommendation
		err := rows.Scan(&rec.ID, &rec.Name, &rec.Description, &rec.MeatCount, 
			&rec.VegetableCount, &rec.IsActive, &rec.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan recommendation"})
			return
		}
		recommendations = append(recommendations, rec)
	}

	c.JSON(http.StatusOK, recommendations)
}

// 获取应季菜品
func (h *Handler) GetSeasonalDishes(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT d.id, d.name, d.description, d.category_id, c.name as category_name,
			   d.price, d.image_url, d.video_url, d.cooking_steps, d.is_seasonal, d.is_active,
			   d.created_at, d.updated_at
		FROM dishes d
		LEFT JOIN categories c ON d.category_id = c.id
		WHERE d.is_seasonal = true AND d.is_active = true
		ORDER BY d.created_at DESC
		LIMIT 10
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seasonal dishes"})
		return
	}
	defer rows.Close()

	var dishes []models.Dish
	for rows.Next() {
		var dish models.Dish
		var categoryName sql.NullString
		
		err := rows.Scan(
			&dish.ID, &dish.Name, &dish.Description, &dish.CategoryID, &categoryName,
			&dish.Price, &dish.ImageURL, &dish.VideoURL, &dish.CookingSteps,
			&dish.IsSeasonal, &dish.IsActive, &dish.CreatedAt, &dish.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan dish"})
			return
		}

		if categoryName.Valid {
			dish.Category = &models.Category{
				ID:   dish.CategoryID,
				Name: categoryName.String,
			}
		}

		dishes = append(dishes, dish)
	}

	c.JSON(http.StatusOK, dishes)
}