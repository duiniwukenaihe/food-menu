package models

import (
    "time"
    "gorm.io/gorm"
)

type Recommendation struct {
    ID         uint           `json:"id" gorm:"primaryKey"`
    UserID     uint           `json:"userId"`
    User       User           `json:"user" gorm:"foreignKey:UserID"`
    ContentID  uint           `json:"contentId"`
    Content    Content        `json:"content" gorm:"foreignKey:ContentID"`
    Score      float64        `json:"score"`
    Reason     string         `json:"reason"`
    IsViewed   bool           `json:"isViewed" gorm:"default:false"`
    CreatedAt  time.Time      `json:"createdAt"`
    UpdatedAt  time.Time      `json:"updatedAt"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// MenuConfig stores configuration for menu recommendations
type MenuConfig struct {
    ID                 uint      `json:"id" gorm:"primaryKey"`
    Name               string    `json:"name" gorm:"type:varchar(100);unique;not null"`
    MeatDishCount      int       `json:"meatDishCount" gorm:"default:1"`
    VegetableDishCount int       `json:"vegetableDishCount" gorm:"default:2"`
    MaxSuggestedDishes int       `json:"maxSuggestedDishes" gorm:"default:6"`
    CreatedAt          time.Time `json:"createdAt"`
    UpdatedAt          time.Time `json:"updatedAt"`
}

type RecommendationResponse struct {
    ID        uint      `json:"id"`
    UserID    uint      `json:"userId"`
    ContentID uint      `json:"contentId"`
    Content   Content   `json:"content"`
    Score     float64   `json:"score"`
    Reason    string    `json:"reason"`
    IsViewed  bool      `json:"isViewed"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

func (r *Recommendation) ToResponse() RecommendationResponse {
    return RecommendationResponse{
        ID:        r.ID,
        UserID:    r.UserID,
        ContentID: r.ContentID,
        Content:   r.Content,
        Score:     r.Score,
        Reason:    r.Reason,
        IsViewed:  r.IsViewed,
        CreatedAt: r.CreatedAt,
        UpdatedAt: r.UpdatedAt,
    }
}
