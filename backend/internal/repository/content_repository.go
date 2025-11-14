package repository

import (
    "context"
    "strings"

    "example.com/app/internal/models"
    "gorm.io/gorm"
)

type ContentFilter struct {
    Search      string
    CategoryID  *uint
    AuthorID    *uint
    IsPublished *bool
    Tags        []string
}

type ContentRepository interface {
    WithContext(ctx context.Context) ContentRepository
    WithTx(tx any) ContentRepository
    DB() any

    Create(content *models.Content) error
    FindByID(id uint) (*models.Content, error)
    Update(content *models.Content) error
    Delete(id uint) error
    List(filter ContentFilter, offset, limit int) ([]models.Content, int64, error)
    IncrementViewCount(id uint) error
}

type contentRepository struct {
    *BaseRepository
}

func NewContentRepository(db *gorm.DB) ContentRepository {
    return &contentRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

func (r *contentRepository) WithContext(ctx context.Context) ContentRepository {
    return &contentRepository{
        BaseRepository: r.BaseRepository.WithContext(ctx),
    }
}

func (r *contentRepository) WithTx(tx any) ContentRepository {
    if txDB, ok := tx.(*gorm.DB); ok {
        return &contentRepository{
            BaseRepository: r.BaseRepository.WithTx(txDB),
        }
    }
    return r
}

func (r *contentRepository) Create(content *models.Content) error {
    return r.DB().Create(content).Error
}

func (r *contentRepository) FindByID(id uint) (*models.Content, error) {
    var content models.Content
    err := r.DB().Preload("Category").Preload("Author").Where("id = ?", id).First(&content).Error
    if err != nil {
        return nil, err
    }
    return &content, nil
}

func (r *contentRepository) Update(content *models.Content) error {
    return r.DB().Save(content).Error
}

func (r *contentRepository) Delete(id uint) error {
    return r.DB().Delete(&models.Content{}, id).Error
}

func (r *contentRepository) List(filter ContentFilter, offset, limit int) ([]models.Content, int64, error) {
    var contents []models.Content
    var total int64

    query := r.DB().Model(&models.Content{}).Preload("Category").Preload("Author")

    if filter.Search != "" {
        searchQuery := strings.ToLower(filter.Search)
        pattern := "%" + searchQuery + "%"
        query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(body) LIKE ?", pattern, pattern, pattern)
    }

    if filter.CategoryID != nil {
        query = query.Where("category_id = ?", *filter.CategoryID)
    }

    if filter.AuthorID != nil {
        query = query.Where("author_id = ?", *filter.AuthorID)
    }

    if filter.IsPublished != nil {
        query = query.Where("is_published = ?", *filter.IsPublished)
    }

    if len(filter.Tags) > 0 {
        for _, tag := range filter.Tags {
            query = query.Where("tags LIKE ?", "%"+tag+"%")
        }
    }

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&contents).Error
    return contents, total, err
}

func (r *contentRepository) IncrementViewCount(id uint) error {
    return r.DB().Model(&models.Content{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
