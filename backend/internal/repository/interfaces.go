package repository

import (
    "context"
    "example.com/app/internal/models"
)

type UserRepository interface {
    WithContext(ctx context.Context) UserRepository
    WithTx(tx any) UserRepository
    DB() any

    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByUsername(username string) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
    List(offset, limit int) ([]models.User, int64, error)
    UpdateLastLogin(id uint) error
}

// Additional repository interfaces can be defined similarly
