//go:build integration
// +build integration

package integration

import (
    "encoding/json"
    "fmt"
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

func setupMenuTestDB() *gorm.DB {
    db, err := database.InitDB()
    if err != nil {
        panic("Failed to connect to test database")
    }

    db.Exec("DELETE FROM dishes")
    db.Exec("DELETE FROM menu_configs")
    db.Exec("DELETE FROM categories")

    return db
}

func setupMenuRouter() *gin.Engine {
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
        v1.GET("/menu/seasonal", api.GetSeasonalMenu)
        v1.GET("/menu/suggested", api.GetSuggestedMenu)
        v1.GET("/menu/config", api.GetMenuConfig)
        v1.GET("/dishes/:id", api.GetDishByID)

        protected := v1.Group("/")
        protected.Use(api.AuthMiddleware())
        {
            admin := protected.Group("/admin")
            admin.Use(api.AdminMiddleware())
            {
                admin.PUT("/menu/config", api.UpdateMenuConfigHandler)
            }
        }
    }

    return r
}

func createTestDishes(db *gorm.DB) {
    category := models.Category{
        Name:     "Main Course",
        IsActive: true,
    }
    db.Create(&category)

    dishes := []models.Dish{
        {
            Name:            "Grilled Chicken",
            Description:     "Delicious grilled chicken",
            Tags:            "meat,protein",
            IsActive:        true,
            CategoryID:      category.ID,
            IsSeasonal:      true,
            AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
            TextSteps:       `[{"step":1,"description":"Prepare chicken"},{"step":2,"description":"Grill"}]`,
        },
        {
            Name:            "Beef Steak",
            Description:     "Premium beef steak",
            Tags:            "meat,protein",
            IsActive:        true,
            CategoryID:      category.ID,
            IsSeasonal:      true,
            AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
            TextSteps:       `[{"step":1,"description":"Season beef"},{"step":2,"description":"Cook"}]`,
        },
        {
            Name:            "Roasted Vegetables",
            Description:     "Mixed roasted vegetables",
            Tags:            "vegetable,veg",
            IsActive:        true,
            CategoryID:      category.ID,
            IsSeasonal:      true,
            AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
            TextSteps:       `[{"step":1,"description":"Prepare veggies"},{"step":2,"description":"Roast"}]`,
        },
        {
            Name:            "Broccoli Salad",
            Description:     "Fresh broccoli salad",
            Tags:            "vegetable,veg,healthy",
            IsActive:        true,
            CategoryID:      category.ID,
            IsSeasonal:      true,
            AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
            TextSteps:       `[{"step":1,"description":"Wash broccoli"},{"step":2,"description":"Mix"}]`,
        },
        {
            Name:            "Fish Fillet",
            Description:     "Fresh fish fillet",
            Tags:            "protein,seafood",
            IsActive:        true,
            CategoryID:      category.ID,
            IsSeasonal:      false,
            TextSteps:       `[{"step":1,"description":"Clean fish"},{"step":2,"description":"Cook"}]`,
        },
        {
            Name:            "Pasta",
            Description:     "Italian pasta",
            Tags:            "carb",
            IsActive:        true,
            CategoryID:      category.ID,
            IsSeasonal:      false,
            TextSteps:       `[{"step":1,"description":"Boil water"},{"step":2,"description":"Cook pasta"}]`,
        },
    }

    for i := range dishes {
        db.Create(&dishes[i])
    }
}

func TestGetSeasonalMenu(t *testing.T) {
    db := setupMenuTestDB()
    api.InitDatabase(db)
    api.InitMenuService()
    router := setupMenuRouter()

    createTestDishes(db)

    config := models.MenuConfig{
        Name:               "default",
        MeatDishCount:      1,
        VegetableDishCount: 2,
        MaxSuggestedDishes: 6,
    }
    db.Create(&config)

    req, _ := http.NewRequest("GET", "/api/v1/menu/seasonal", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var paginatedResponse models.PaginatedResponse
    err := json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
    assert.NoError(t, err)
    assert.True(t, paginatedResponse.Success)
    assert.Greater(t, paginatedResponse.Total, int64(0))
    assert.Greater(t, len(paginatedResponse.Data.([]interface{})), 0)
}

func TestGetSuggestedMenu(t *testing.T) {
    db := setupMenuTestDB()
    api.InitDatabase(db)
    api.InitMenuService()
    router := setupMenuRouter()

    createTestDishes(db)

    config := models.MenuConfig{
        Name:               "default",
        MeatDishCount:      1,
        VegetableDishCount: 2,
        MaxSuggestedDishes: 6,
    }
    db.Create(&config)

    req, _ := http.NewRequest("GET", "/api/v1/menu/suggested", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var paginatedResponse models.PaginatedResponse
    err := json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
    assert.NoError(t, err)
    assert.True(t, paginatedResponse.Success)
    assert.Greater(t, paginatedResponse.Total, int64(0))
    assert.LessOrEqual(t, paginatedResponse.Total, int64(6))
}

func TestSuggestedMenuLimit(t *testing.T) {
    db := setupMenuTestDB()
    api.InitDatabase(db)
    api.InitMenuService()
    router := setupMenuRouter()

    category := models.Category{Name: "Main Course", IsActive: true}
    db.Create(&category)

    for i := 0; i < 10; i++ {
        dish := models.Dish{
            Name:       "Dish " + string(rune('A'+i)),
            IsActive:   true,
            CategoryID: category.ID,
        }
        db.Create(&dish)
    }

    config := models.MenuConfig{
        Name:               "default",
        MeatDishCount:      1,
        VegetableDishCount: 2,
        MaxSuggestedDishes: 3,
    }
    db.Create(&config)

    req, _ := http.NewRequest("GET", "/api/v1/menu/suggested", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var paginatedResponse models.PaginatedResponse
    err := json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
    assert.NoError(t, err)
    assert.LessOrEqual(t, paginatedResponse.Total, int64(3))
}

func TestSeasonalMenuPagination(t *testing.T) {
    db := setupMenuTestDB()
    api.InitDatabase(db)
    api.InitMenuService()
    router := setupMenuRouter()

    createTestDishes(db)

    config := models.MenuConfig{
        Name:               "default",
        MeatDishCount:      1,
        VegetableDishCount: 2,
        MaxSuggestedDishes: 6,
    }
    db.Create(&config)

    req, _ := http.NewRequest("GET", "/api/v1/menu/seasonal?page=1&limit=2", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var paginatedResponse models.PaginatedResponse
    err := json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
    assert.NoError(t, err)
    assert.Equal(t, 1, paginatedResponse.Page)
    assert.Equal(t, 2, paginatedResponse.Limit)
    assert.LessOrEqual(t, len(paginatedResponse.Data.([]interface{})), 2)
}

func TestSuggestedMenuPagination(t *testing.T) {
    db := setupMenuTestDB()
    api.InitDatabase(db)
    api.InitMenuService()
    router := setupMenuRouter()

    createTestDishes(db)

    config := models.MenuConfig{
        Name:               "default",
        MeatDishCount:      1,
        VegetableDishCount: 2,
        MaxSuggestedDishes: 6,
    }
    db.Create(&config)

    req, _ := http.NewRequest("GET", "/api/v1/menu/suggested?page=1&limit=2", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var paginatedResponse models.PaginatedResponse
    err := json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
    assert.NoError(t, err)
    assert.Equal(t, 1, paginatedResponse.Page)
    assert.Equal(t, 2, paginatedResponse.Limit)
    assert.LessOrEqual(t, len(paginatedResponse.Data.([]interface{})), 2)
}

func TestGetMenuConfig(t *testing.T) {
    db := setupMenuTestDB()
    api.InitDatabase(db)
    api.InitMenuService()
    router := setupMenuRouter()

    config := models.MenuConfig{
        Name:               "default",
        MeatDishCount:      2,
        VegetableDishCount: 3,
        MaxSuggestedDishes: 8,
    }
    db.Create(&config)

    req, _ := http.NewRequest("GET", "/api/v1/menu/config", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var successResponse models.SuccessResponse
    err := json.Unmarshal(w.Body.Bytes(), &successResponse)
    assert.NoError(t, err)
    assert.True(t, successResponse.Success)

    configData := successResponse.Data.(map[string]interface{})
    assert.Equal(t, float64(2), configData["meatDishCount"])
    assert.Equal(t, float64(3), configData["vegetableDishCount"])
    assert.Equal(t, float64(8), configData["maxSuggestedDishes"])
}

func TestDishDetail(t *testing.T) {
    db := setupMenuTestDB()
    api.InitDatabase(db)
    api.InitMenuService()
    router := setupMenuRouter()

    category := models.Category{Name: "Main Course", IsActive: true}
    db.Create(&category)

    dish := models.Dish{
        Name:        "Test Dish",
        Description: "Test description",
        CategoryID:  category.ID,
        IsActive:    true,
        TextSteps:   `[{"step":1,"description":"Step 1"}]`,
    }
    db.Create(&dish)

    dishURL := fmt.Sprintf("/api/v1/dishes/%d", dish.ID)
    req, _ := http.NewRequest("GET", dishURL, nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var successResponse models.SuccessResponse
    err := json.Unmarshal(w.Body.Bytes(), &successResponse)
    assert.NoError(t, err)
    assert.True(t, successResponse.Success)
}
