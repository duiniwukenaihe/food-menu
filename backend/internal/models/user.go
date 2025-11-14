package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"unique;not null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"-" gorm:"not null"`
    FirstName string         `json:"firstName"`
    LastName  string         `json:"lastName"`
    Role      string         `json:"role" gorm:"default:'user'"`
    Avatar    string         `json:"avatar"`
    IsActive  bool           `json:"isActive" gorm:"default:true"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    Role      string    `json:"role"`
    Avatar    string    `json:"avatar"`
    IsActive  bool      `json:"isActive"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        FirstName: u.FirstName,
        LastName:  u.LastName,
        Role:      u.Role,
        Avatar:    u.Avatar,
        IsActive:  u.IsActive,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username  string `json:"username" binding:"required,min=3,max=50"`
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8,strong_password"`
    FirstName string `json:"firstName" binding:"required"`
    LastName  string `json:"lastName" binding:"required"`
}

type UpdateProfileRequest struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Avatar    string `json:"avatar"`
}

type AuthResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}
