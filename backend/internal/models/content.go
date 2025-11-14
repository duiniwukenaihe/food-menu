package models

import (
    "time"
    "gorm.io/gorm"
)

type Content struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Title       string         `json:"title" gorm:"type:varchar(255);not null;index"`
    Description string         `json:"description" gorm:"type:text"`
    Body        string         `json:"body" gorm:"type:text"`
    CategoryID  uint           `json:"categoryId" gorm:"not null;index"`
    Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
    AuthorID    uint           `json:"authorId" gorm:"not null;index"`
    Author      User           `json:"author" gorm:"foreignKey:AuthorID"`
    Tags        string         `json:"tags" gorm:"type:varchar(500)"`
    ImageURL    string         `json:"imageUrl" gorm:"type:varchar(500)"`
    IsPublished bool           `json:"isPublished" gorm:"default:false;index"`
    ViewCount   int            `json:"viewCount" gorm:"default:0;index"`
    UpdatedBy   *uint          `json:"updatedBy" gorm:"index"`
    CreatedAt   time.Time      `json:"createdAt" gorm:"index"`
    UpdatedAt   time.Time      `json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ContentResponse struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Body        string    `json:"body"`
    CategoryID  uint      `json:"categoryId"`
    Category    Category  `json:"category"`
    AuthorID    uint      `json:"authorId"`
    Author      User      `json:"author"`
    Tags        string    `json:"tags"`
    ImageURL    string    `json:"imageUrl"`
    IsPublished bool      `json:"isPublished"`
    ViewCount   int       `json:"viewCount"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

func (c *Content) ToResponse() ContentResponse {
    return ContentResponse{
        ID:          c.ID,
        Title:       c.Title,
        Description: c.Description,
        Body:        c.Body,
        CategoryID:  c.CategoryID,
        Category:    c.Category,
        AuthorID:    c.AuthorID,
        Author:      c.Author,
        Tags:        c.Tags,
        ImageURL:    c.ImageURL,
        IsPublished: c.IsPublished,
        ViewCount:   c.ViewCount,
        CreatedAt:   c.CreatedAt,
        UpdatedAt:   c.UpdatedAt,
    }
}

type CreateContentRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=200"`
    Description string `json:"description"`
    Body        string `json:"body" binding:"required"`
    CategoryID  uint   `json:"categoryId" binding:"required"`
    Tags        string `json:"tags"`
    ImageURL    string `json:"imageUrl"`
    IsPublished bool   `json:"isPublished"`
}

type UpdateContentRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Body        string `json:"body"`
    CategoryID  uint   `json:"categoryId"`
    Tags        string `json:"tags"`
    ImageURL    string `json:"imageUrl"`
    IsPublished *bool  `json:"isPublished"`
}
