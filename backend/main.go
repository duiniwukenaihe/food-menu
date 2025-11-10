package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "example.com/app/internal/api"
    "example.com/app/internal/database"
    "example.com/app/internal/models"
)

// @title Full-Stack App API
// @version 1.0
// @description A modern web application with authentication, content management, and recommendations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // Initialize database
    db, err := database.InitDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize API with database
    api.InitDatabase(db)

    // Initialize storage service
    if err := api.InitStorageService(); err != nil {
        log.Fatal("Failed to initialize storage service:", err)
    }

    // Auto-migrate the schema
    err = db.AutoMigrate(&models.User{}, &models.Content{}, &models.Category{}, &models.Recommendation{}, &models.Dish{}, &models.Media{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Set Gin mode
    if os.Getenv("ENVIRONMENT") == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Initialize router
    r := gin.Default()

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

    // API routes
    v1 := r.Group("/api/v1")
    {
        // Public routes
        v1.POST("/auth/register", api.Register)
        v1.POST("/auth/login", api.Login)
        v1.GET("/content", api.GetContent)
        v1.GET("/content/:id", api.GetContentByID)
        v1.GET("/categories", api.GetCategories)
        v1.GET("/recommendations", api.GetRecommendations)
        v1.GET("/dishes", api.GetDishes)
        v1.GET("/dishes/:id", api.GetDishByID)
        
        // Media serving (public)
        v1.GET("/media/*filepath", api.GetMedia)

        // Protected routes
        protected := v1.Group("/")
        protected.Use(api.AuthMiddleware())
        {
            protected.GET("/auth/profile", api.GetProfile)
            protected.PUT("/auth/profile", api.UpdateProfile)

            // Admin routes
            admin := protected.Group("/admin")
            admin.Use(api.AdminMiddleware())
            {
                admin.GET("/users", api.GetUsers)
                admin.POST("/users", api.CreateUser)
                admin.PUT("/users/:id", api.UpdateUser)
                admin.DELETE("/users/:id", api.DeleteUser)

                admin.GET("/content", api.AdminGetContent)
                admin.POST("/content", api.CreateContent)
                admin.PUT("/content/:id", api.UpdateContent)
                admin.DELETE("/content/:id", api.DeleteContent)

                admin.GET("/categories", api.AdminGetCategories)
                admin.POST("/categories", api.CreateCategory)
                admin.PUT("/categories/:id", api.UpdateCategory)
                admin.DELETE("/categories/:id", api.DeleteCategory)
                
                // Dish management
                admin.GET("/dishes", api.AdminGetDishes)
                admin.POST("/dishes", api.CreateDish)
                admin.PUT("/dishes/:id", api.UpdateDish)
                admin.DELETE("/dishes/:id", api.DeleteDish)
                
                // Media upload
                admin.POST("/media/upload-url", api.GetUploadURL)
                admin.POST("/media/upload", api.UploadFile)
                admin.DELETE("/media", api.DeleteMedia)
            }
        }
    }

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
