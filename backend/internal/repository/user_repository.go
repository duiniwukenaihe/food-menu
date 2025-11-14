package repository

import (
	"context"
	"time"

	"example.com/app/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *userRepository) WithContext(ctx context.Context) UserRepository {
	return &userRepository{
		BaseRepository: r.BaseRepository.WithContext(ctx),
	}
}

func (r *userRepository) WithTx(tx any) UserRepository {
	if txDB, ok := tx.(*gorm.DB); ok {
		return &userRepository{
			BaseRepository: r.BaseRepository.WithTx(txDB),
		}
	}
	return r
}

func (r *userRepository) Create(user *models.User) error {
	return r.DB().Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.DB().Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.DB().Delete(&models.User{}, id).Error
}

func (r *userRepository) List(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := r.DB().Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.DB().Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	now := time.Now()
	return r.DB().Model(&models.User{}).Where("id = ?", id).Update("last_login", now).Error
}
