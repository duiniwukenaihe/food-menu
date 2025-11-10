package models

import (
    "time"
)

// Dish represents a dish in the system
type Dish struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"type:varchar(255);not null"`
    Description string    `json:"description" gorm:"type:text"`
    Tags        string    `json:"tags" gorm:"type:varchar(500)"` // Comma-separated tags
    IsActive    bool      `json:"isActive" gorm:"default:true"`
    
    // Category association
    CategoryID  uint     `json:"categoryId"`
    Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
    
    // Recipe/cooking details
    TextSteps   string    `json:"textSteps" gorm:"type:text"` // JSON array of step objects
    
    // Seasonal configuration
    IsSeasonal       bool   `json:"isSeasonal" gorm:"default:false"`
    AvailableMonths  string `json:"availableMonths" gorm:"type:varchar(100)"` // e.g., "1,2,3,11,12" for Jan-Mar and Nov-Dec
    SeasonalNote     string `json:"seasonalNote" gorm:"type:varchar(500)"`
    
    // Media URLs
    ImageURL         string `json:"imageUrl" gorm:"type:varchar(500)"`
    ThumbnailURL     string `json:"thumbnailUrl" gorm:"type:varchar(500)"`
    GalleryURLs      string `json:"galleryUrls" gorm:"type:text"` // JSON array of image URLs
    
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// CreateDishRequest represents the request body for creating a dish
type CreateDishRequest struct {
    Name             string `json:"name" binding:"required,min=1,max=255"`
    Description      string `json:"description"`
    Tags             string `json:"tags"`
    IsActive         bool   `json:"isActive"`
    CategoryID       uint   `json:"categoryId"`
    TextSteps        string `json:"textSteps"`
    IsSeasonal       bool   `json:"isSeasonal"`
    AvailableMonths  string `json:"availableMonths"`
    SeasonalNote     string `json:"seasonalNote"`
    ImageURL         string `json:"imageUrl"`
    ThumbnailURL     string `json:"thumbnailUrl"`
    GalleryURLs      string `json:"galleryUrls"`
}

// UpdateDishRequest represents the request body for updating a dish
type UpdateDishRequest struct {
    Name             *string `json:"name" binding:"omitempty,min=1,max=255"`
    Description      *string `json:"description"`
    Tags             *string `json:"tags"`
    IsActive         *bool   `json:"isActive"`
    CategoryID       *uint   `json:"categoryId"`
    TextSteps        *string `json:"textSteps"`
    IsSeasonal       *bool   `json:"isSeasonal"`
    AvailableMonths  *string `json:"availableMonths"`
    SeasonalNote     *string `json:"seasonalNote"`
    ImageURL         *string `json:"imageUrl"`
    ThumbnailURL     *string `json:"thumbnailUrl"`
    GalleryURLs      *string `json:"galleryUrls"`
}

// DishQueryParams represents query parameters for filtering dishes
type DishQueryParams struct {
    PaginationQuery
    Search     string `form:"search"`
    IsSeasonal *bool  `form:"isSeasonal"`
    IsActive   *bool  `form:"isActive"`
    Tags       string `form:"tags"`
    CategoryID uint   `form:"categoryId"`
}
