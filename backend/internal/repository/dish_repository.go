package repository

import (
    "context"
    "strings"

    "example.com/app/internal/models"
    "gorm.io/gorm"
)

type DishFilter struct {
    Search     string
    CategoryID *uint
    IsSeasonal *bool
    IsActive   *bool
    Tags       []string
}

type DishRepository interface {
    WithContext(ctx context.Context) DishRepository
    WithTx(tx any) DishRepository
    DB() any

    Create(dish *models.Dish) error
    FindByID(id uint) (*models.Dish, error)
    Update(dish *models.Dish) error
    Delete(id uint) error
    List(filter DishFilter, offset, limit int) ([]models.Dish, int64, error)
}

type dishRepository struct {
    *BaseRepository
}

func NewDishRepository(db *gorm.DB) DishRepository {
    return &dishRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

func (r *dishRepository) WithContext(ctx context.Context) DishRepository {
    return &dishRepository{
        BaseRepository: r.BaseRepository.WithContext(ctx),
    }
}

func (r *dishRepository) WithTx(tx any) DishRepository {
    if txDB, ok := tx.(*gorm.DB); ok {
        return &dishRepository{
            BaseRepository: r.BaseRepository.WithTx(txDB),
        }
    }
    return r
}

func (r *dishRepository) Create(dish *models.Dish) error {
    return r.DB().Create(dish).Error
}

func (r *dishRepository) FindByID(id uint) (*models.Dish, error) {
    var dish models.Dish
    err := r.DB().Preload("Category").Where("id = ?", id).First(&dish).Error
    if err != nil {
        return nil, err
    }
    return &dish, nil
}

func (r *dishRepository) Update(dish *models.Dish) error {
    return r.DB().Save(dish).Error
}

func (r *dishRepository) Delete(id uint) error {
    return r.DB().Delete(&models.Dish{}, id).Error
}

func (r *dishRepository) List(filter DishFilter, offset, limit int) ([]models.Dish, int64, error) {
    var dishes []models.Dish
    var total int64

    query := r.DB().Model(&models.Dish{}).Preload("Category")

    if filter.Search != "" {
        searchQuery := strings.ToLower(filter.Search)
        pattern := "%" + searchQuery + "%"
        query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?", pattern, pattern, pattern)
    }

    if filter.CategoryID != nil {
        query = query.Where("category_id = ?", *filter.CategoryID)
    }

    if filter.IsSeasonal != nil {
        query = query.Where("is_seasonal = ?", *filter.IsSeasonal)
    }

    if filter.IsActive != nil {
        query = query.Where("is_active = ?", *filter.IsActive)
    }

    if len(filter.Tags) > 0 {
        for _, tag := range filter.Tags {
            tag = strings.TrimSpace(tag)
            if tag != "" {
                query = query.Where("tags LIKE ?", "%"+tag+"%")
            }
        }
    }

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&dishes).Error
    return dishes, total, err
}
