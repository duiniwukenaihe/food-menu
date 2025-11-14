package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"unique;not null;index"`
    Email     string         `json:"email" gorm:"unique;not null;index"`
    Password  string         `json:"-" gorm:"column:password_hash;not null"`
    FirstName string         `json:"firstName" gorm:"type:varchar(100)"`
    LastName  string         `json:"lastName" gorm:"type:varchar(100)"`
    Role      string         `json:"role" gorm:"type:varchar(50);default:'user';index"`
    Avatar    string         `json:"avatar" gorm:"type:varchar(500)"`
    IsActive  bool           `json:"isActive" gorm:"default:true;index"`
    LastLogin *time.Time     `json:"lastLogin"`
    CreatedAt time.Time      `json:"createdAt" gorm:"index"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserResponse struct {
    ID        uint       `json:"id"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
    FirstName string     `json:"firstName"`
    LastName  string     `json:"lastName"`
    Role      string     `json:"role"`
    Avatar    string     `json:"avatar"`
    IsActive  bool       `json:"isActive"`
    LastLogin *time.Time `json:"lastLogin"`
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
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
        LastLogin: u.LastLogin,
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
    Password  string `json:"password" binding:"required,min=6"`
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
