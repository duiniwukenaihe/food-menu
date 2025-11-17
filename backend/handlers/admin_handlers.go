package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"food-ordering/models"

	"github.com/gin-gonic/gin"
)

// 获取用户列表（管理员）
func (h *Handler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.DefaultQuery("search", "")
	
	offset := (page - 1) * limit

	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM users
	`
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += ` WHERE username ILIKE $1 OR email ILIKE $2`
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex = 3
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
			return
		}
		users = append(users, user)
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM users"
	countArgs := []interface{}{}
	if search != "" {
		countQuery += " WHERE username ILIKE $1 OR email ILIKE $2"
		countArgs = append(countArgs, "%"+search+"%", "%"+search+"%")
	}

	var total int
	h.db.QueryRow(countQuery, countArgs...).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// 创建菜品（管理员）
func (h *Handler) CreateDish(c *gin.Context) {
	var req models.CreateDishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证分类是否存在
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", req.CategoryID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var dishID int
	err = h.db.QueryRow(`
		INSERT INTO dishes (name, description, category_id, price, image_url, video_url, 
			cooking_steps, is_seasonal, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true, NOW(), NOW())
		RETURNING id
	`, req.Name, req.Description, req.CategoryID, req.Price, req.ImageURL, 
		req.VideoURL, req.CookingSteps, req.IsSeasonal).Scan(&dishID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dish"})
		return
	}

	// 返回创建的菜品
	dish, err := h.getDishByID(dishID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dish created but failed to fetch details"})
		return
	}

	c.JSON(http.StatusCreated, dish)
}

// 更新菜品（管理员）
func (h *Handler) UpdateDish(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dish ID"})
		return
	}

	var req models.UpdateDishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建动态更新查询
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		updates = append(updates, "name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, "description = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.CategoryID != nil {
		// 验证分类是否存在
		var exists bool
		err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", *req.CategoryID).Scan(&exists)
		if err != nil || !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		updates = append(updates, "category_id = $"+strconv.Itoa(argIndex))
		args = append(args, *req.CategoryID)
		argIndex++
	}
	if req.Price != nil {
		updates = append(updates, "price = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Price)
		argIndex++
	}
	if req.ImageURL != nil {
		updates = append(updates, "image_url = $"+strconv.Itoa(argIndex))
		args = append(args, *req.ImageURL)
		argIndex++
	}
	if req.VideoURL != nil {
		updates = append(updates, "video_url = $"+strconv.Itoa(argIndex))
		args = append(args, *req.VideoURL)
		argIndex++
	}
	if req.CookingSteps != nil {
		updates = append(updates, "cooking_steps = $"+strconv.Itoa(argIndex))
		args = append(args, *req.CookingSteps)
		argIndex++
	}
	if req.IsSeasonal != nil {
		updates = append(updates, "is_seasonal = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsSeasonal)
		argIndex++
	}
	if req.IsActive != nil {
		updates = append(updates, "is_active = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	updates = append(updates, "updated_at = NOW()")

	query := "UPDATE dishes SET " + join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)
	args = append(args, id)

	result, err := h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dish"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	// 返回更新后的菜品
	dish, err := h.getDishByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dish updated but failed to fetch details"})
		return
	}

	c.JSON(http.StatusOK, dish)
}

// 删除菜品（管理员）
func (h *Handler) DeleteDish(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dish ID"})
		return
	}

	// 软删除：设置为不活跃
	result, err := h.db.Exec("UPDATE dishes SET is_active = false, updated_at = NOW() WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete dish"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
}

// 创建分类（管理员）
func (h *Handler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var categoryID int
	err := h.db.QueryRow(`
		INSERT INTO categories (name, description, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id
	`, req.Name, req.Description).Scan(&categoryID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	// 返回创建的分类
	category, err := h.getCategoryByID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Category created but failed to fetch details"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// 更新分类（管理员）
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		updates = append(updates, "name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, "description = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query := "UPDATE categories SET " + join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)
	args = append(args, id)

	result, err := h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// 返回更新后的分类
	category, err := h.getCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Category updated but failed to fetch details"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// 删除分类（管理员）
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// 检查是否有菜品使用此分类
	var count int
	err = h.db.QueryRow("SELECT COUNT(*) FROM dishes WHERE category_id = $1 AND is_active = true", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check category usage"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with active dishes"})
		return
	}

	result, err := h.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// 获取配置（管理员）
func (h *Handler) GetConfig(c *gin.Context) {
	rows, err := h.db.Query("SELECT id, config_key, config_value, description, updated_at FROM system_config ORDER BY config_key")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch config"})
		return
	}
	defer rows.Close()

	var configs []models.SystemConfig
	for rows.Next() {
		var config models.SystemConfig
		err := rows.Scan(&config.ID, &config.ConfigKey, &config.ConfigValue, &config.Description, &config.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan config"})
			return
		}
		configs = append(configs, config)
	}

	c.JSON(http.StatusOK, configs)
}

// 更新配置（管理员）
func (h *Handler) UpdateConfig(c *gin.Context) {
	var configs map[string]string
	if err := c.ShouldBindJSON(&configs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback()

	for key, value := range configs {
		_, err := tx.Exec(`
			UPDATE system_config 
			SET config_value = $1, updated_at = NOW() 
			WHERE config_key = $2
		`, value, key)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config: " + key})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Config updated successfully"})
}

// 辅助方法
func (h *Handler) getDishByID(id int) (*models.Dish, error) {
	var dish models.Dish
	var categoryName sql.NullString

	err := h.db.QueryRow(`
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
		return nil, err
	}

	if categoryName.Valid {
		dish.Category = &models.Category{
			ID:   dish.CategoryID,
			Name: categoryName.String,
		}
	}

	return &dish, nil
}

func (h *Handler) getCategoryByID(id int) (*models.Category, error) {
	var category models.Category
	err := h.db.QueryRow("SELECT id, name, description, created_at FROM categories WHERE id = $1", id).
		Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
	
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}