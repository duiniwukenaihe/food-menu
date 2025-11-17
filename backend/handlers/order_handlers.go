package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"food-ordering/models"

	"github.com/gin-gonic/gin"
)

// 创建订单
func (h *Handler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order must contain at least one item"})
		return
	}

	// 开始事务
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback()

	// 计算总金额并验证菜品
	var totalAmount float64
	for _, item := range req.Items {
		var price float64
		err := tx.QueryRow("SELECT price FROM dishes WHERE id = $1 AND is_active = true", item.DishID).Scan(&price)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Dish %d not found", item.DishID)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dish price"})
			return
		}
		totalAmount += price * float64(item.Quantity)
	}

	// 创建订单
	var orderID int
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, total_amount, status) 
		VALUES ($1, $2, 'pending') 
		RETURNING id
	`, userID, totalAmount).Scan(&orderID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// 创建订单明细
	for _, item := range req.Items {
		var price float64
		tx.QueryRow("SELECT price FROM dishes WHERE id = $1", item.DishID).Scan(&price)
		
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, dish_id, quantity, price) 
			VALUES ($1, $2, $3, $4)
		`, orderID, item.DishID, item.Quantity, price)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// 返回创建的订单
	order, err := h.getOrderWithItems(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order created but failed to fetch details"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// 获取用户订单
func (h *Handler) GetOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	offset := (page - 1) * limit

	rows, err := h.db.Query(`
		SELECT id, user_id, total_amount, status, created_at, updated_at
		FROM orders 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}
		orders = append(orders, order)
	}

	// 获取总数
	var total int
	h.db.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = $1", userID).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// 获取订单详情（辅助方法）
func (h *Handler) getOrderWithItems(orderID int) (*models.Order, error) {
	var order models.Order
	
	// 获取订单基本信息
	err := h.db.QueryRow(`
		SELECT id, user_id, total_amount, status, created_at, updated_at
		FROM orders WHERE id = $1
	`, orderID).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	
	if err != nil {
		return nil, err
	}

	// 获取订单明细
	rows, err := h.db.Query(`
		SELECT oi.id, oi.order_id, oi.dish_id, oi.quantity, oi.price, oi.created_at,
			   d.name, d.image_url
		FROM order_items oi
		LEFT JOIN dishes d ON oi.dish_id = d.id
		WHERE oi.order_id = $1
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		var dishName sql.NullString
		var dishImageURL sql.NullString
		
		err := rows.Scan(&item.ID, &item.OrderID, &item.DishID, &item.Quantity, 
			&item.Price, &item.CreatedAt, &dishName, &dishImageURL)
		if err != nil {
			return nil, err
		}

		if dishName.Valid {
			item.Dish = &models.Dish{
				ID:       item.DishID,
				Name:     dishName.String,
				ImageURL: dishImageURL.String,
			}
		}

		order.Items = append(order.Items, item)
	}

	return &order, nil
}

// 添加到收藏
func (h *Handler) AddToFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
	dishID, err := strconv.Atoi(c.Param("dishId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dish ID"})
		return
	}

	// 检查菜品是否存在
	var exists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM dishes WHERE id = $1 AND is_active = true)", dishID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	// 检查是否已经收藏
	var alreadyExists bool
	err = h.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM user_favorites WHERE user_id = $1 AND dish_id = $2)
	`, userID, dishID).Scan(&alreadyExists)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if alreadyExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Dish already in favorites"})
		return
	}

	// 添加到收藏
	_, err = h.db.Exec(`
		INSERT INTO user_favorites (user_id, dish_id) 
		VALUES ($1, $2)
	`, userID, dishID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites"})
}

// 从收藏中移除
func (h *Handler) RemoveFromFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
	dishID, err := strconv.Atoi(c.Param("dishId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dish ID"})
		return
	}

	result, err := h.db.Exec(`
		DELETE FROM user_favorites 
		WHERE user_id = $1 AND dish_id = $2
	`, userID, dishID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from favorites"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

// 获取用户收藏
func (h *Handler) GetFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	offset := (page - 1) * limit

	rows, err := h.db.Query(`
		SELECT uf.id, uf.user_id, uf.dish_id, uf.created_at,
			   d.name, d.description, d.price, d.image_url, d.is_seasonal,
			   c.name as category_name
		FROM user_favorites uf
		LEFT JOIN dishes d ON uf.dish_id = d.id
		LEFT JOIN categories c ON d.category_id = c.id
		WHERE uf.user_id = $1 AND d.is_active = true
		ORDER BY uf.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}
	defer rows.Close()

	var favorites []models.UserFavorite
	for rows.Next() {
		var fav models.UserFavorite
		var dishName, dishDesc, dishImageURL, categoryName sql.NullString
		var dishPrice sql.NullFloat64
		var isSeasonal sql.NullBool
		
		err := rows.Scan(
			&fav.ID, &fav.UserID, &fav.DishID, &fav.CreatedAt,
			&dishName, &dishDesc, &dishPrice, &dishImageURL, &isSeasonal, &categoryName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan favorite"})
			return
		}

		if dishName.Valid {
			fav.Dish = &models.Dish{
				ID:          fav.DishID,
				Name:        dishName.String,
				Description: dishDesc.String,
				Price:       dishPrice.Float64,
				ImageURL:    dishImageURL.String,
				IsSeasonal:  isSeasonal.Bool,
			}
			if categoryName.Valid {
				fav.Dish.Category = &models.Category{Name: categoryName.String}
			}
		}

		favorites = append(favorites, fav)
	}

	// 获取总数
	var total int
	h.db.QueryRow(`
		SELECT COUNT(*) FROM user_favorites uf
		LEFT JOIN dishes d ON uf.dish_id = d.id
		WHERE uf.user_id = $1 AND d.is_active = true
	`, userID).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"favorites": favorites,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}