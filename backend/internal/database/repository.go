package database

import (
	"context"
	"example.com/app/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	WithContext(ctx context.Context) Repository
	WithTx(tx *gorm.DB) Repository
	DB() *gorm.DB
}

type UserRepository interface {
	Repository
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(offset, limit int) ([]models.User, int64, error)
	UpdateLastLogin(id uint) error
}

type DishRepository interface {
	Repository
	Create(dish *models.Dish) error
	FindByID(id uint) (*models.Dish, error)
	Update(dish *models.Dish) error
	Delete(id uint) error
	List(filters map[string]interface{}, offset, limit int) ([]models.Dish, int64, error)
	Search(query string, offset, limit int) ([]models.Dish, int64, error)
	FindByCategory(categoryID uint, offset, limit int) ([]models.Dish, int64, error)
	FindSeasonal(offset, limit int) ([]models.Dish, int64, error)
}

type DishCategoryRepository interface {
	Repository
	Create(category *models.DishCategory) error
	FindByID(id uint) (*models.DishCategory, error)
	FindByName(name string) (*models.DishCategory, error)
	Update(category *models.DishCategory) error
	Delete(id uint) error
	List(offset, limit int) ([]models.DishCategory, int64, error)
}

type DishPairingRepository interface {
	Repository
	Create(pairing *models.DishPairing) error
	FindByID(id uint) (*models.DishPairing, error)
	FindByDishID(dishID uint) ([]models.DishPairing, error)
	Delete(id uint) error
	DeleteByDishID(dishID uint) error
}

type CategoryRepository interface {
	Repository
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindByName(name string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	List(offset, limit int) ([]models.Category, int64, error)
}

type ContentRepository interface {
	Repository
	Create(content *models.Content) error
	FindByID(id uint) (*models.Content, error)
	Update(content *models.Content) error
	Delete(id uint) error
	List(filters map[string]interface{}, offset, limit int) ([]models.Content, int64, error)
	IncrementViewCount(id uint) error
	FindByAuthor(authorID uint, offset, limit int) ([]models.Content, int64, error)
	FindByCategory(categoryID uint, offset, limit int) ([]models.Content, int64, error)
}

type MediaRepository interface {
	Repository
	Create(media *models.Media) error
	FindByID(id uint) (*models.Media, error)
	FindByKey(key string) (*models.Media, error)
	Update(media *models.Media) error
	Delete(id uint) error
	DeleteByKey(key string) error
	FindByEntity(entityType string, entityID uint) ([]models.Media, error)
}

type RecommendationRepository interface {
	Repository
	Create(recommendation *models.Recommendation) error
	FindByID(id uint) (*models.Recommendation, error)
	Update(recommendation *models.Recommendation) error
	Delete(id uint) error
	FindByUserID(userID uint, offset, limit int) ([]models.Recommendation, int64, error)
	MarkAsViewed(id uint) error
	DeleteOldRecommendations(userID uint, keepCount int) error
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
