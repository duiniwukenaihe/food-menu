package repository

// Examples of repository usage patterns
// These are documentation examples, not functional code

/*
Example 1: Basic CRUD operations with context

	ctx := context.Background()
	userRepo := repository.NewUserRepository(db)

	// Create
	user := &models.User{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}
	err := userRepo.WithContext(ctx).Create(user)

	// Read
	foundUser, err := userRepo.WithContext(ctx).FindByUsername("john_doe")

	// Update
	foundUser.FirstName = "John"
	err = userRepo.WithContext(ctx).Update(foundUser)

	// Delete (soft delete)
	err = userRepo.WithContext(ctx).Delete(foundUser.ID)

Example 2: Using repository filters

	ctx := context.Background()
	dishRepo := repository.NewDishRepository(db)

	filter := repository.DishFilter{
		Search:     "pizza",
		IsActive:   &trueVal,
		IsSeasonal: &falseVal,
		Tags:       []string{"italian", "cheese"},
	}

	dishes, total, err := dishRepo.WithContext(ctx).List(filter, 0, 10)

Example 3: Using the Store pattern

	store := repository.NewStore(db)
	ctx := context.Background()

	// Access different repositories through store
	user, err := store.Users().WithContext(ctx).FindByID(1)
	dishes, total, err := store.Dishes().WithContext(ctx).List(filter, 0, 10)

Example 4: Transaction handling

	txManager := repository.NewTransactionManager(db)
	ctx := context.Background()

	err := txManager.WithTransaction(ctx, func(txCtx context.Context, txStore *Store) error {
		// All operations in this function use the same transaction
		user := &models.User{Username: "newuser", Email: "new@example.com"}
		if err := txStore.Users().WithContext(txCtx).Create(user); err != nil {
			return err // Transaction will be rolled back
		}

		dish := &models.Dish{Name: "New Dish", CreatedBy: &user.ID}
		if err := txStore.Dishes().WithContext(txCtx).Create(dish); err != nil {
			return err // Transaction will be rolled back
		}

		return nil // Transaction will be committed
	})

Example 5: Content repository with filters

	contentRepo := store.Contents().WithContext(ctx)

	filter := ContentFilter{
		Search:      "recipe",
		IsPublished: &trueVal,
		CategoryID:  &catID,
		Tags:        []string{"healthy", "quick"},
	}

	contents, total, err := contentRepo.List(filter, 0, 20)

	// Increment view count
	err = contentRepo.IncrementViewCount(contentID)

Example 6: Using repositories in API handlers

	func GetDishByID(c *gin.Context) {
		ctx := c.Request.Context()
		dishRepo := repository.NewDishRepository(db)

		id, _ := strconv.Atoi(c.Param("id"))
		dish, err := dishRepo.WithContext(ctx).FindByID(uint(id))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
			return
		}

		c.JSON(http.StatusOK, dish)
	}

Example 7: Complex queries with multiple filters

	ctx := context.Background()
	dishRepo := repository.NewDishRepository(db)

	active := true
	seasonal := true
	categoryID := uint(5)

	filter := repository.DishFilter{
		Search:     "summer",
		CategoryID: &categoryID,
		IsActive:   &active,
		IsSeasonal: &seasonal,
		Tags:       []string{"fresh", "light"},
	}

	// Get first page
	dishes, total, err := dishRepo.WithContext(ctx).List(filter, 0, 10)

	// Calculate pagination
	totalPages := (total + 9) / 10
	hasNextPage := total > 10

Example 8: Transaction with error handling

	txManager := repository.NewTransactionManager(db)

	err := txManager.WithTransaction(ctx, func(txCtx context.Context, txStore *Store) error {
		// Find user
		user, err := txStore.Users().WithContext(txCtx).FindByID(userID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		// Update user
		user.LastLogin = &now
		if err := txStore.Users().WithContext(txCtx).Update(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		// Create activity log (if you have such a model)
		// ...

		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	}

Example 9: Batch operations within transaction

	err := txManager.WithTransaction(ctx, func(txCtx context.Context, txStore *Store) error {
		// Create multiple dishes
		dishes := []*models.Dish{
			{Name: "Dish 1", CreatedBy: &userID},
			{Name: "Dish 2", CreatedBy: &userID},
			{Name: "Dish 3", CreatedBy: &userID},
		}

		for _, dish := range dishes {
			if err := txStore.Dishes().WithContext(txCtx).Create(dish); err != nil {
				return err // Rolls back all creations
			}
		}

		return nil
	})

Example 10: Repository in service layer

	type DishService struct {
		dishRepo repository.DishRepository
		userRepo repository.UserRepository
	}

	func NewDishService(store *repository.Store) *DishService {
		return &DishService{
			dishRepo: store.Dishes(),
			userRepo: store.Users(),
		}
	}

	func (s *DishService) CreateDish(ctx context.Context, dish *models.Dish, userID uint) error {
		// Verify user exists
		user, err := s.userRepo.WithContext(ctx).FindByID(userID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		// Set audit fields
		dish.CreatedBy = &userID
		now := time.Now()
		dish.CreatedAt = now
		dish.UpdatedAt = now

		// Create dish
		return s.dishRepo.WithContext(ctx).Create(dish)
	}
*/
