package models

import (
    "time"
    "gorm.io/gorm"
)

type Recommendation struct {
    ID         uint           `json:"id" gorm:"primaryKey"`
    UserID     uint           `json:"userId" gorm:"not null;index"`
    User       User           `json:"user" gorm:"foreignKey:UserID"`
    ContentID  *uint          `json:"contentId" gorm:"index"`
    Content    *Content       `json:"content,omitempty" gorm:"foreignKey:ContentID"`
    DishID     *uint          `json:"dishId" gorm:"index"`
    Dish       *Dish          `json:"dish,omitempty" gorm:"foreignKey:DishID"`
    EntityType string         `json:"entityType" gorm:"type:varchar(50);not null;index"` // "content" or "dish"
    Score      float64        `json:"score" gorm:"index"`
    Reason     string         `json:"reason" gorm:"type:text"`
    Algorithm  string         `json:"algorithm" gorm:"type:varchar(50)"` // e.g., "popularity", "collaborative", "content-based"
    IsViewed   bool           `json:"isViewed" gorm:"default:false;index"`
    ViewedAt   *time.Time     `json:"viewedAt"`
    CreatedAt  time.Time      `json:"createdAt" gorm:"index"`
    UpdatedAt  time.Time      `json:"updatedAt"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type RecommendationResponse struct {
    ID         uint        `json:"id"`
    UserID     uint        `json:"userId"`
    ContentID  *uint       `json:"contentId"`
    Content    *Content    `json:"content,omitempty"`
    DishID     *uint       `json:"dishId"`
    Dish       *Dish       `json:"dish,omitempty"`
    EntityType string      `json:"entityType"`
    Score      float64     `json:"score"`
    Reason     string      `json:"reason"`
    Algorithm  string      `json:"algorithm"`
    IsViewed   bool        `json:"isViewed"`
    ViewedAt   *time.Time  `json:"viewedAt"`
    CreatedAt  time.Time   `json:"createdAt"`
    UpdatedAt  time.Time   `json:"updatedAt"`
}

func (r *Recommendation) ToResponse() RecommendationResponse {
    return RecommendationResponse{
        ID:         r.ID,
        UserID:     r.UserID,
        ContentID:  r.ContentID,
        Content:    r.Content,
        DishID:     r.DishID,
        Dish:       r.Dish,
        EntityType: r.EntityType,
        Score:      r.Score,
        Reason:     r.Reason,
        Algorithm:  r.Algorithm,
        IsViewed:   r.IsViewed,
        ViewedAt:   r.ViewedAt,
        CreatedAt:  r.CreatedAt,
        UpdatedAt:  r.UpdatedAt,
    }
}
