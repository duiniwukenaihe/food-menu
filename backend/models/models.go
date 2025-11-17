package models

import (
	"database/sql"
	"time"
)

// 用户模型
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 菜品分类模型
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// 菜品模型
type Dish struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CategoryID    int       `json:"category_id"`
	Category      *Category `json:"category,omitempty"`
	Price         float64   `json:"price"`
	ImageURL      string    `json:"image_url"`
	VideoURL      string    `json:"video_url"`
	CookingSteps  string    `json:"cooking_steps"`
	IsSeasonal    bool      `json:"is_seasonal"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// 菜品营养信息
type DishNutrition struct {
	ID           int     `json:"id"`
	DishID       int     `json:"dish_id"`
	Calories     int     `json:"calories"`
	Protein      float64 `json:"protein"`
	Fat          float64 `json:"fat"`
	Carbohydrates float64 `json:"carbohydrates"`
	Fiber        float64 `json:"fiber"`
	CreatedAt    time.Time `json:"created_at"`
}

// 订单模型
type Order struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Items       []OrderItem `json:"items,omitempty"`
}

// 订单明细
type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	DishID    int       `json:"dish_id"`
	Dish      *Dish     `json:"dish,omitempty"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// 推荐配置
type Recommendation struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	MeatCount        int       `json:"meat_count"`
	VegetableCount   int       `json:"vegetable_count"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	Dishes           []Dish    `json:"dishes,omitempty"`
}

// 用户收藏
type UserFavorite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	DishID    int       `json:"dish_id"`
	Dish      *Dish     `json:"dish,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// 系统配置
type SystemConfig struct {
	ID          int       `json:"id"`
	ConfigKey   string    `json:"config_key"`
	ConfigValue string    `json:"config_value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User  `json:"user"`
}

// 创建订单请求
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required"`
}

type CreateOrderItemRequest struct {
	DishID   int `json:"dish_id" binding:"required"`
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// 创建菜品请求
type CreateDishRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	CategoryID   int     `json:"category_id" binding:"required"`
	Price        float64 `json:"price" binding:"required,min=0"`
	ImageURL     string  `json:"image_url"`
	VideoURL     string  `json:"video_url"`
	CookingSteps string  `json:"cooking_steps"`
	IsSeasonal   bool    `json:"is_seasonal"`
}

// 更新菜品请求
type UpdateDishRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	CategoryID   *int     `json:"category_id"`
	Price        *float64 `json:"price"`
	ImageURL     *string  `json:"image_url"`
	VideoURL     *string  `json:"video_url"`
	CookingSteps *string  `json:"cooking_steps"`
	IsSeasonal   *bool    `json:"is_seasonal"`
	IsActive     *bool    `json:"is_active"`
}

// 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// 更新分类请求
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// 数据库自动迁移
func AutoMigrate(db *sql.DB) error {
	// 这里应该执行schema.sql文件，但为了简化，我们假设表已经存在
	return nil
}