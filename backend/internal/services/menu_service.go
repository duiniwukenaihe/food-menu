package services

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.com/app/internal/models"
	"gorm.io/gorm"
)

type MenuService struct {
	db *gorm.DB
}

func NewMenuService(database *gorm.DB) *MenuService {
	return &MenuService{
		db: database,
	}
}

// GetSeasonalRecommendation generates a seasonal menu recommendation with configurable combinations
func (m *MenuService) GetSeasonalRecommendation(page, limit int) ([]models.Dish, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	// Get configuration from database, with fallback defaults
	meatCount := 1
	vegCount := 2

	var config models.MenuConfig
	if err := m.db.Where("name = ?", "default").First(&config).Error; err == nil {
		meatCount = config.MeatDishCount
		vegCount = config.VegetableDishCount
	}

	// Get currently available seasonal dishes
	var seasonalDishes []models.Dish
	currentMonth := time.Now().Month()

	if err := m.db.Where("is_seasonal = ? AND is_active = ?", true, true).Find(&seasonalDishes).Error; err != nil {
		return nil, 0, err
	}

	// Filter by availability month
	var availableDishes []models.Dish
	for _, dish := range seasonalDishes {
		if m.isDishAvailableInMonth(dish.AvailableMonths, currentMonth) {
			availableDishes = append(availableDishes, dish)
		}
	}

	// Select random combos respecting the configuration
	selectedDishes := m.selectRandomCombination(availableDishes, meatCount, vegCount)

	// Apply pagination to selection
	total := int64(len(selectedDishes))
	offset := (page - 1) * limit

	if offset >= len(selectedDishes) {
		return []models.Dish{}, total, nil
	}

	endIdx := offset + limit
	if endIdx > len(selectedDishes) {
		endIdx = len(selectedDishes)
	}

	return selectedDishes[offset:endIdx], total, nil
}

// GetSuggestedDishes generates "guess you like" randomized dish recommendations
func (m *MenuService) GetSuggestedDishes(userID *uint, page, limit int) ([]models.Dish, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	// Get configuration from database, with fallback defaults
	maxSuggested := 6

	var config models.MenuConfig
	if err := m.db.Where("name = ?", "default").First(&config).Error; err == nil {
		maxSuggested = config.MaxSuggestedDishes
	}

	// Get all active dishes
	var dishes []models.Dish
	if err := m.db.Where("is_active = ?", true).Find(&dishes).Error; err != nil {
		return nil, 0, err
	}

	if len(dishes) == 0 {
		return []models.Dish{}, 0, nil
	}

	// Randomize selection
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(dishes), func(i, j int) {
		dishes[i], dishes[j] = dishes[j], dishes[i]
	})

	// Limit to maxSuggested dishes
	if len(dishes) > maxSuggested {
		dishes = dishes[:maxSuggested]
	}

	// Apply pagination
	total := int64(len(dishes))
	offset := (page - 1) * limit

	if offset >= len(dishes) {
		return []models.Dish{}, total, nil
	}

	endIdx := offset + limit
	if endIdx > len(dishes) {
		endIdx = len(dishes)
	}

	return dishes[offset:endIdx], total, nil
}

// selectRandomCombination selects random combination of meat and vegetable dishes
func (m *MenuService) selectRandomCombination(dishes []models.Dish, meatCount, vegCount int) []models.Dish {
	var meatDishes []models.Dish
	var vegDishes []models.Dish
	var otherDishes []models.Dish

	// Categorize dishes based on tags
	for _, dish := range dishes {
		tags := strings.ToLower(dish.Tags)
		if strings.Contains(tags, "meat") || strings.Contains(tags, "protein") {
			meatDishes = append(meatDishes, dish)
		} else if strings.Contains(tags, "vegetable") || strings.Contains(tags, "veg") {
			vegDishes = append(vegDishes, dish)
		} else {
			otherDishes = append(otherDishes, dish)
		}
	}

	// Shuffle and select
	selected := []models.Dish{}

	// Select meat dishes
	rand.Shuffle(len(meatDishes), func(i, j int) {
		meatDishes[i], meatDishes[j] = meatDishes[j], meatDishes[i]
	})
	for i := 0; i < meatCount && i < len(meatDishes); i++ {
		selected = append(selected, meatDishes[i])
	}

	// Select vegetable dishes
	rand.Shuffle(len(vegDishes), func(i, j int) {
		vegDishes[i], vegDishes[j] = vegDishes[j], vegDishes[i]
	})
	for i := 0; i < vegCount && i < len(vegDishes); i++ {
		selected = append(selected, vegDishes[i])
	}

	// If we don't have enough dishes, fill with others
	if len(selected) < meatCount+vegCount {
		needed := (meatCount + vegCount) - len(selected)
		rand.Shuffle(len(otherDishes), func(i, j int) {
			otherDishes[i], otherDishes[j] = otherDishes[j], otherDishes[i]
		})
		for i := 0; i < needed && i < len(otherDishes); i++ {
			selected = append(selected, otherDishes[i])
		}
	}

	return selected
}

// isDishAvailableInMonth checks if a dish is available in the given month
func (m *MenuService) isDishAvailableInMonth(availableMonths string, month time.Month) bool {
	if availableMonths == "" {
		return true
	}

	months := strings.Split(availableMonths, ",")
	targetMonth := strconv.Itoa(int(month))

	for _, m := range months {
		if strings.TrimSpace(m) == targetMonth {
			return true
		}
	}

	return false
}

// GetMenuConfig retrieves the default menu configuration
func (m *MenuService) GetMenuConfig() (*models.MenuConfig, error) {
	var config models.MenuConfig
	if err := m.db.Where("name = ?", "default").First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default configuration if not found
			return &models.MenuConfig{
				Name:               "default",
				MeatDishCount:      1,
				VegetableDishCount: 2,
				MaxSuggestedDishes: 6,
			}, nil
		}
		return nil, err
	}
	return &config, nil
}

// UpdateMenuConfig updates the default menu configuration
func (m *MenuService) UpdateMenuConfig(config *models.MenuConfig) error {
	config.Name = "default"
	return m.db.Save(config).Error
}
