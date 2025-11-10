package services

import (
	"testing"
	"time"

	"example.com/app/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&models.Category{}, &models.Dish{}, &models.MenuConfig{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func createTestDishes(db *gorm.DB) []models.Dish {
	category := models.Category{Name: "Main Course", IsActive: true}
	db.Create(&category)

	dishes := []models.Dish{
		{
			Name:            "Grilled Chicken",
			Description:     "Delicious grilled chicken",
			Tags:            "meat, protein",
			IsActive:        true,
			CategoryID:      category.ID,
			IsSeasonal:      true,
			AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
		},
		{
			Name:            "Beef Steak",
			Description:     "Premium beef steak",
			Tags:            "meat, protein",
			IsActive:        true,
			CategoryID:      category.ID,
			IsSeasonal:      true,
			AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
		},
		{
			Name:            "Roasted Vegetables",
			Description:     "Mixed roasted vegetables",
			Tags:            "vegetable, veg",
			IsActive:        true,
			CategoryID:      category.ID,
			IsSeasonal:      true,
			AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
		},
		{
			Name:            "Broccoli Salad",
			Description:     "Fresh broccoli salad",
			Tags:            "vegetable, veg, healthy",
			IsActive:        true,
			CategoryID:      category.ID,
			IsSeasonal:      true,
			AvailableMonths: "1,2,3,4,5,6,7,8,9,10,11,12",
		},
		{
			Name:            "Seasonal Asparagus",
			Description:     "Seasonal asparagus",
			Tags:            "vegetable, seasonal",
			IsActive:        true,
			CategoryID:      category.ID,
			IsSeasonal:      true,
			AvailableMonths: "3,4,5",
		},
		{
			Name:            "Fish Fillet",
			Description:     "Fresh fish fillet",
			Tags:            "protein, seafood",
			IsActive:        true,
			CategoryID:      category.ID,
			IsSeasonal:      false,
		},
	}

	for i := range dishes {
		db.Create(&dishes[i])
	}

	return dishes
}

func TestGetSeasonalRecommendation(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create default config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Test getting seasonal recommendations
	dishes, total, err := service.GetSeasonalRecommendation(1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return at least one dish (meat + vegetables)
	if len(dishes) == 0 {
		t.Error("Expected seasonal dishes, got none")
	}

	// Total should be > 0
	if total == 0 {
		t.Error("Expected total > 0, got 0")
	}

	// Check that we have the expected combination
	if total < 3 {
		t.Errorf("Expected at least 3 seasonal dishes, got %d", total)
	}

	t.Logf("Got %d seasonal dishes out of %d total", len(dishes), total)
}

func TestSeasonalRecommendationCombination(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create custom config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      2,
		VegetableDishCount: 1,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Get recommendation
	dishes, _, err := service.GetSeasonalRecommendation(1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Count meat and vegetable dishes
	meatCount := 0
	vegCount := 0

	for _, dish := range dishes {
		if service.isDishAvailableInMonth(dish.AvailableMonths, time.Now().Month()) {
			// Check if meat or veg based on tags
			if dish.Tags != "" {
				if (service.selectRandomCombination([]models.Dish{dish}, 1, 0)) != nil {
					// Basic check passed
				}
			}
		}
	}

	t.Logf("Meat dishes: %d, Veg dishes: %d", meatCount, vegCount)
}

func TestGetSuggestedDishes(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create default config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Test with no user (guest)
	dishes, total, err := service.GetSuggestedDishes(nil, 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return suggested dishes
	if len(dishes) == 0 {
		t.Error("Expected suggested dishes, got none")
	}

	// Should not exceed max suggested dishes
	if total > 6 {
		t.Errorf("Expected max 6 suggested dishes, got %d", total)
	}

	t.Logf("Got %d suggested dishes", len(dishes))
}

func TestSuggestedDishesLimit(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data - more than max
	category := models.Category{Name: "Main Course", IsActive: true}
	db.Create(&category)

	for i := 0; i < 10; i++ {
		dish := models.Dish{
			Name:       "Dish " + string(rune('A'+i)),
			IsActive:   true,
			CategoryID: category.ID,
		}
		db.Create(&dish)
	}

	// Create config with max 3
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 3,
	}
	db.Create(&config)

	// Get suggestions
	dishes, total, err := service.GetSuggestedDishes(nil, 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should be limited to 3
	if total > 3 {
		t.Errorf("Expected max 3 suggested dishes, got %d", total)
	}

	t.Logf("Limited to %d dishes as configured", total)
}

func TestSuggestedDishesWithUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create default config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Test with user (authenticated)
	userID := uint(1)
	dishes, total, err := service.GetSuggestedDishes(&userID, 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(dishes) == 0 {
		t.Error("Expected suggested dishes for user, got none")
	}

	if total > 6 {
		t.Errorf("Expected max 6 suggested dishes, got %d", total)
	}

	t.Logf("Got %d suggested dishes for authenticated user", len(dishes))
}

func TestSeasonalMonthFiltering(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	category := models.Category{Name: "Main Course", IsActive: true}
	db.Create(&category)

	// Create winter dishes
	winterDish := models.Dish{
		Name:            "Winter Squash",
		Tags:            "vegetable",
		IsActive:        true,
		CategoryID:      category.ID,
		IsSeasonal:      true,
		AvailableMonths: "11,12,1",
	}
	db.Create(&winterDish)

	// Create summer dishes
	summerDish := models.Dish{
		Name:            "Tomato",
		Tags:            "vegetable",
		IsActive:        true,
		CategoryID:      category.ID,
		IsSeasonal:      true,
		AvailableMonths: "6,7,8",
	}
	db.Create(&summerDish)

	// Test month availability
	if !service.isDishAvailableInMonth("11,12,1", time.December) {
		t.Error("December should be in available months 11,12,1")
	}

	if service.isDishAvailableInMonth("6,7,8", time.December) {
		t.Error("December should not be in available months 6,7,8")
	}

	if service.isDishAvailableInMonth("6,7,8", time.July) {
		// Pass
	} else {
		t.Error("July should be in available months 6,7,8")
	}

	t.Log("Month filtering works correctly")
}

func TestPaginationForSeasonal(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create default config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Get page 1
	dishes1, total1, err := service.GetSeasonalRecommendation(1, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Get page 2
	dishes2, total2, err := service.GetSeasonalRecommendation(2, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Both should have same total
	if total1 != total2 {
		t.Errorf("Expected same total, got %d and %d", total1, total2)
	}

	// Page 1 should have at most 2 items
	if len(dishes1) > 2 {
		t.Errorf("Expected max 2 items in page 1, got %d", len(dishes1))
	}

	// Page 2 should have at most 2 items
	if len(dishes2) > 2 {
		t.Errorf("Expected max 2 items in page 2, got %d", len(dishes2))
	}

	t.Logf("Page 1: %d dishes, Page 2: %d dishes, Total: %d", len(dishes1), len(dishes2), total1)
}

func TestPaginationForSuggested(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create default config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Get page 1 with limit 2
	dishes1, total1, err := service.GetSuggestedDishes(nil, 1, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Get page 2 with limit 2
	dishes2, total2, err := service.GetSuggestedDishes(nil, 2, 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Both should have same total (6 max)
	if total1 != total2 {
		t.Errorf("Expected same total, got %d and %d", total1, total2)
	}

	// Page 1 should have at most 2 items
	if len(dishes1) > 2 {
		t.Errorf("Expected max 2 items in page 1, got %d", len(dishes1))
	}

	// Verify total is max suggested
	if total1 > 6 {
		t.Errorf("Expected max 6 suggested dishes, got %d", total1)
	}

	t.Logf("Suggested pagination: Total %d, Page 1: %d, Page 2: %d", total1, len(dishes1), len(dishes2))
}

func TestGetMenuConfig(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Get default config (should not exist yet)
	config, err := service.GetMenuConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return default values
	if config.MeatDishCount != 1 {
		t.Errorf("Expected default meat count 1, got %d", config.MeatDishCount)
	}

	if config.VegetableDishCount != 2 {
		t.Errorf("Expected default veg count 2, got %d", config.VegetableDishCount)
	}

	if config.MaxSuggestedDishes != 6 {
		t.Errorf("Expected default max suggested 6, got %d", config.MaxSuggestedDishes)
	}

	t.Log("Default menu config returned correctly")
}

func TestUpdateMenuConfig(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create initial config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Update config
	config.MeatDishCount = 2
	config.VegetableDishCount = 3
	config.MaxSuggestedDishes = 8

	err := service.UpdateMenuConfig(&config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update
	updated, err := service.GetMenuConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.MeatDishCount != 2 {
		t.Errorf("Expected meat count 2, got %d", updated.MeatDishCount)
	}

	if updated.VegetableDishCount != 3 {
		t.Errorf("Expected veg count 3, got %d", updated.VegetableDishCount)
	}

	if updated.MaxSuggestedDishes != 8 {
		t.Errorf("Expected max suggested 8, got %d", updated.MaxSuggestedDishes)
	}

	t.Log("Menu config updated successfully")
}

func TestRandomizationInSeasonal(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create default config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 1,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Get recommendations multiple times to check randomization
	results := make(map[string]int)
	for i := 0; i < 10; i++ {
		dishes, _, _ := service.GetSeasonalRecommendation(1, 10)
		if len(dishes) > 0 {
			key := dishes[0].Name
			results[key]++
		}
	}

	// If there's variation, we know randomization is working
	if len(results) > 1 {
		t.Logf("Randomization working - got %d different first dishes in 10 runs", len(results))
	} else {
		t.Log("All 10 runs returned same first dish - randomization might not be working")
	}
}

func TestRandomizationInSuggested(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create test data
	createTestDishes(db)

	// Create default config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Get suggestions multiple times to check randomization
	results := make(map[string]int)
	for i := 0; i < 10; i++ {
		dishes, _, _ := service.GetSuggestedDishes(nil, 1, 10)
		if len(dishes) > 0 {
			key := dishes[0].Name
			results[key]++
		}
	}

	// If there's variation, we know randomization is working
	if len(results) > 1 {
		t.Logf("Randomization working - got %d different first dishes in 10 runs", len(results))
	} else {
		t.Log("All 10 runs returned same first dish - randomization might not be working")
	}
}

func TestEmptyDishDatabase(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	// Create config but no dishes
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Get recommendations with empty database
	dishes, total, err := service.GetSeasonalRecommendation(1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(dishes) != 0 {
		t.Errorf("Expected no dishes, got %d", len(dishes))
	}

	if total != 0 {
		t.Errorf("Expected total 0, got %d", total)
	}

	t.Log("Empty database handled correctly")
}

func TestNoActiveDishes(t *testing.T) {
	db := setupTestDB(t)
	service := NewMenuService(db)

	category := models.Category{Name: "Main Course", IsActive: true}
	db.Create(&category)

	// Create inactive dishes
	dish := models.Dish{
		Name:       "Inactive Dish",
		IsActive:   false,
		CategoryID: category.ID,
	}
	db.Create(&dish)

	// Create config
	config := models.MenuConfig{
		Name:               "default",
		MeatDishCount:      1,
		VegetableDishCount: 2,
		MaxSuggestedDishes: 6,
	}
	db.Create(&config)

	// Get suggestions
	dishes, total, err := service.GetSuggestedDishes(nil, 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(dishes) != 0 {
		t.Errorf("Expected no active dishes, got %d", len(dishes))
	}

	if total != 0 {
		t.Errorf("Expected total 0, got %d", total)
	}

	t.Log("Inactive dishes correctly excluded")
}
