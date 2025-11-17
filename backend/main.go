package main

import (
	"food-ordering/config"
	"food-ordering/database"
	"food-ordering/handlers"
	"food-ordering/middleware"
	"food-ordering/models"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// 自动迁移数据库表
	if err := models.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化处理器
	handler := handlers.NewHandler(db, cfg)

	// 设置Gin路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开路由
		public := api.Group("/")
		{
			public.POST("/login", handler.Login)
			public.GET("/dishes", handler.GetDishes)
			public.GET("/dishes/:id", handler.GetDish)
			public.GET("/categories", handler.GetCategories)
			public.GET("/recommendations", handler.GetRecommendations)
			public.GET("/seasonal-dishes", handler.GetSeasonalDishes)
		}

		// 需要认证的路由
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

		// 管理员路由
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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}