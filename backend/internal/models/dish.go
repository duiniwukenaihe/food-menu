package models

import (
    "time"
    "gorm.io/gorm"
)

// Dish represents a dish in the system
type Dish struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"type:varchar(255);not null;index"`
    Description string         `json:"description" gorm:"type:text"`
    CategoryID  *uint          `json:"categoryId" gorm:"index"`
    Category    *DishCategory  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
    Tags        string         `json:"tags" gorm:"type:varchar(500)"` // Comma-separated tags
    IsActive    bool           `json:"isActive" gorm:"default:true;index"`
    
    // Seasonal configuration
    IsSeasonal       bool   `json:"isSeasonal" gorm:"default:false;index"`
    AvailableMonths  string `json:"availableMonths" gorm:"type:varchar(100)"` // e.g., "1,2,3,11,12" for Jan-Mar and Nov-Dec
    SeasonalNote     string `json:"seasonalNote" gorm:"type:varchar(500)"`
    
    // Recipe details
    Ingredients  string `json:"ingredients" gorm:"type:text"` // JSON array of ingredients
    CookingSteps string `json:"cookingSteps" gorm:"type:text"` // Cooking instructions
    PrepTime     int    `json:"prepTime"` // Preparation time in minutes
    CookTime     int    `json:"cookTime"` // Cooking time in minutes
    Servings     int    `json:"servings" gorm:"default:1"` // Number of servings
    Difficulty   string `json:"difficulty" gorm:"type:varchar(50)"` // e.g., "easy", "medium", "hard"
    
    // Media URLs
    ImageURL         string `json:"imageUrl" gorm:"type:varchar(500)"`
    ThumbnailURL     string `json:"thumbnailUrl" gorm:"type:varchar(500)"`
    GalleryURLs      string `json:"galleryUrls" gorm:"type:text"` // JSON array of image URLs
    VideoURL         string `json:"videoUrl" gorm:"type:varchar(500)"` // Video URL
    
    // Admin audit fields
    CreatedBy *uint         `json:"createdBy" gorm:"index"`
    UpdatedBy *uint         `json:"updatedBy" gorm:"index"`
    CreatedAt time.Time     `json:"createdAt" gorm:"index"`
    UpdatedAt time.Time     `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// DishCategory represents categories for dishes
type DishCategory struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"type:varchar(100);not null;unique;index"`
    Description string         `json:"description" gorm:"type:text"`
    Icon        string         `json:"icon" gorm:"type:varchar(100)"`
    SortOrder   int            `json:"sortOrder" gorm:"default:0"`
    IsActive    bool           `json:"isActive" gorm:"default:true"`
    CreatedAt   time.Time      `json:"createdAt"`
    UpdatedAt   time.Time      `json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// DishPairing represents pairing rules between dishes
type DishPairing struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    DishID       uint      `json:"dishId" gorm:"not null;index"`
    Dish         Dish      `json:"dish" gorm:"foreignKey:DishID"`
    PairedDishID uint      `json:"pairedDishId" gorm:"not null;index"`
    PairedDish   Dish      `json:"pairedDish" gorm:"foreignKey:PairedDishID"`
    PairingType  string    `json:"pairingType" gorm:"type:varchar(50)"` // e.g., "complements", "alternative", "beverage"
    Score        float64   `json:"score" gorm:"default:1.0"` // Pairing strength score
    CreatedAt    time.Time `json:"createdAt"`
}

// CreateDishRequest represents the request body for creating a dish
type CreateDishRequest struct {
    Name             string `json:"name" binding:"required,min=1,max=255"`
    Description      string `json:"description"`
    CategoryID       *uint  `json:"categoryId"`
    Tags             string `json:"tags"`
    IsActive         bool   `json:"isActive"`
    IsSeasonal       bool   `json:"isSeasonal"`
    AvailableMonths  string `json:"availableMonths"`
    SeasonalNote     string `json:"seasonalNote"`
    Ingredients      string `json:"ingredients"`
    CookingSteps     string `json:"cookingSteps"`
    PrepTime         int    `json:"prepTime"`
    CookTime         int    `json:"cookTime"`
    Servings         int    `json:"servings"`
    Difficulty       string `json:"difficulty"`
    ImageURL         string `json:"imageUrl"`
    ThumbnailURL     string `json:"thumbnailUrl"`
    GalleryURLs      string `json:"galleryUrls"`
    VideoURL         string `json:"videoUrl"`
}

// UpdateDishRequest represents the request body for updating a dish
type UpdateDishRequest struct {
    Name             *string `json:"name" binding:"omitempty,min=1,max=255"`
    Description      *string `json:"description"`
    CategoryID       *uint   `json:"categoryId"`
    Tags             *string `json:"tags"`
    IsActive         *bool   `json:"isActive"`
    IsSeasonal       *bool   `json:"isSeasonal"`
    AvailableMonths  *string `json:"availableMonths"`
    SeasonalNote     *string `json:"seasonalNote"`
    Ingredients      *string `json:"ingredients"`
    CookingSteps     *string `json:"cookingSteps"`
    PrepTime         *int    `json:"prepTime"`
    CookTime         *int    `json:"cookTime"`
    Servings         *int    `json:"servings"`
    Difficulty       *string `json:"difficulty"`
    ImageURL         *string `json:"imageUrl"`
    ThumbnailURL     *string `json:"thumbnailUrl"`
    GalleryURLs      *string `json:"galleryUrls"`
    VideoURL         *string `json:"videoUrl"`
}

// DishQueryParams represents query parameters for filtering dishes
type DishQueryParams struct {
    PaginationQuery
    Search     string `form:"search"`
    IsSeasonal *bool  `form:"isSeasonal"`
    IsActive   *bool  `form:"isActive"`
    Tags       string `form:"tags"`
}
