package api

import (
    "math/rand"
    "net/http"
    "strconv"
    "time"
    "example.com/app/internal/models"
    "github.com/gin-gonic/gin"
)

// GetRecommendations godoc
// @Summary Get user recommendations
// @Description Get personalized content recommendations for the authenticated user
// @Tags recommendations
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of recommendations" default(5)
// @Success 200 {object} models.SuccessResponse
// @Router /recommendations [get]
func GetRecommendations(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, models.ErrorResponse{
            Success: false,
            Message: "User not authenticated",
        })
        return
    }

    // Simple recommendation algorithm based on popular content and randomization
    var contents []models.Content
    
    // Get popular content (high view count)
    if err := db.Preload("Category").Preload("Author").
        Where("is_published = ?", true).
        Order("view_count DESC").
        Limit(20).
        Find(&contents).Error; err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Success: false,
            Message: "Failed to fetch recommendations",
            Error:   err.Error(),
        })
        return
    }

    // Randomize and limit the results
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(contents), func(i, j int) {
        contents[i], contents[j] = contents[j], contents[i]
    })

    limit := 5
    if l := c.Query("limit"); l != "" {
        if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 20 {
            limit = parsedLimit
        }
    }

    if len(contents) > limit {
        contents = contents[:limit]
    }

    // Convert to response format
    var contentResponses []models.ContentResponse
    for _, content := range contents {
        contentResponses = append(contentResponses, content.ToResponse())
    }

    // Create recommendation records for tracking
    for _, content := range contents {
        recommendation := models.Recommendation{
            UserID:    userID.(uint),
            ContentID: content.ID,
            Score:     rand.Float64() * 100,
            Reason:    "Based on popular content",
            IsViewed:  false,
        }
        db.Create(&recommendation)
    }

    c.JSON(http.StatusOK, models.SuccessResponse{
        Success: true,
        Data:    contentResponses,
    })
}
