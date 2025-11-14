package models

import (
    "time"
    "gorm.io/gorm"
)

type Category struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"type:varchar(100);unique;not null;index"`
    Description string         `json:"description" gorm:"type:text"`
    Color       string         `json:"color" gorm:"type:varchar(50)"`
    Icon        string         `json:"icon" gorm:"type:varchar(100)"`
    IsActive    bool           `json:"isActive" gorm:"default:true;index"`
    SortOrder   int            `json:"sortOrder" gorm:"default:0"`
    CreatedBy   *uint          `json:"createdBy" gorm:"index"`
    UpdatedBy   *uint          `json:"updatedBy" gorm:"index"`
    CreatedAt   time.Time      `json:"createdAt" gorm:"index"`
    UpdatedAt   time.Time      `json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CategoryResponse struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Color       string    `json:"color"`
    Icon        string    `json:"icon"`
    SortOrder   int       `json:"sortOrder"`
    IsActive    bool      `json:"isActive"`
    CreatedBy   *uint     `json:"createdBy"`
    UpdatedBy   *uint     `json:"updatedBy"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

func (c *Category) ToResponse() CategoryResponse {
    return CategoryResponse{
        ID:          c.ID,
        Name:        c.Name,
        Description: c.Description,
        Color:       c.Color,
        Icon:        c.Icon,
        SortOrder:   c.SortOrder,
        IsActive:    c.IsActive,
        CreatedBy:   c.CreatedBy,
        UpdatedBy:   c.UpdatedBy,
        CreatedAt:   c.CreatedAt,
        UpdatedAt:   c.UpdatedAt,
    }
}

type CreateCategoryRequest struct {
    Name        string `json:"name" binding:"required,min=1,max=100"`
    Description string `json:"description"`
    Color       string `json:"color"`
    Icon        string `json:"icon"`
}

type UpdateCategoryRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Color       string `json:"color"`
    Icon        string `json:"icon"`
    IsActive    *bool  `json:"isActive"`
}
