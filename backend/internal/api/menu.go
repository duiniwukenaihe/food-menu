package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "example.com/app/internal/models"
    "example.com/app/internal/services"
)

var menuService *services.MenuService

func InitMenuService() {
    menuService = services.NewMenuService(db)
}

// GetSeasonalMenu godoc
// @Summary Get seasonal menu recommendations
// @Description Get randomly selected seasonal dishes based on configurable combination rules
// @Tags menu
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /menu/seasonal [get]
func GetSeasonalMenu(c *gin.Context) {
    var params models.DishQueryParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid query parameters",
            Error:   err.Error(),
        })
        return
    }

    if params.Limit <= 0 {
        params.Limit = 10
    }
    if params.Page <= 0 {
        params.Page = 1
    }

    dishes, total, err := menuService.GetSeasonalRecommendation(params.Page, params.Limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to fetch seasonal menu",
            Error:   err.Error(),
        })
        return
    }

    // Ensure dishes have category data populated
    for i := range dishes {
        db.Model(&dishes[i]).Association("Category").Find(&dishes[i].Category)
    }

    c.JSON(http.StatusOK, models.PaginatedResponse{
        Success: true,
        Data:    dishes,
        Total:   total,
        Page:    params.Page,
        Limit:   params.Limit,
    })
}

// GetSuggestedMenu godoc
// @Summary Get suggested menu recommendations
// @Description Get randomized "guess you like" dish recommendations limited to configured maximum
// @Tags menu
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /menu/suggested [get]
func GetSuggestedMenu(c *gin.Context) {
    var params models.DishQueryParams
    if err := c.ShouldBindQuery(&params); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid query parameters",
            Error:   err.Error(),
        })
        return
    }

    if params.Limit <= 0 {
        params.Limit = 10
    }
    if params.Page <= 0 {
        params.Page = 1
    }

    // Check if user is authenticated
    var userID *uint
    if u, exists := c.Get("userID"); exists {
        uid := u.(uint)
        userID = &uid
    }

    dishes, total, err := menuService.GetSuggestedDishes(userID, params.Page, params.Limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to fetch suggested menu",
            Error:   err.Error(),
        })
        return
    }

    // Ensure dishes have category data populated
    for i := range dishes {
        db.Model(&dishes[i]).Association("Category").Find(&dishes[i].Category)
    }

    c.JSON(http.StatusOK, models.PaginatedResponse{
        Success: true,
        Data:    dishes,
        Total:   total,
        Page:    params.Page,
        Limit:   params.Limit,
    })
}

// GetMenuConfig godoc
// @Summary Get menu configuration
// @Description Get the current menu recommendation configuration
// @Tags menu
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /menu/config [get]
func GetMenuConfig(c *gin.Context) {
    config, err := menuService.GetMenuConfig()
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to fetch menu config",
            Error:   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, models.SuccessResponse{
        Success: true,
        Data:    config,
    })
}

// UpdateMenuConfig godoc
// @Summary Update menu configuration
// @Description Update the menu recommendation configuration (admin only)
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param config body models.MenuConfig true "Menu config data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/menu/config [put]
func UpdateMenuConfigHandler(c *gin.Context) {
    var req models.MenuConfig
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Success: false,
            Message: "Invalid request body",
            Error:   err.Error(),
        })
        return
    }

    if err := menuService.UpdateMenuConfig(&req); err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to update menu config",
            Error:   err.Error(),
        })
        return
    }

    config, _ := menuService.GetMenuConfig()

    c.JSON(http.StatusOK, models.SuccessResponse{
        Success: true,
        Message: "Menu config updated successfully",
        Data:    config,
    })
}
