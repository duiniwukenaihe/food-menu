package api

import (
    "context"

    "example.com/app/internal/repository"
)

var repoStore *repository.Store

// InitRepository initializes the repository store used by API handlers.
func InitRepository(store *repository.Store) {
    repoStore = store
}

func dishRepository(ctx context.Context) repository.DishRepository {
    if repoStore != nil {
        return repoStore.Dishes().WithContext(ctx)
    }
    if db != nil {
        return repository.NewDishRepository(db).WithContext(ctx)
    }
    return nil
}

func userRepository(ctx context.Context) repository.UserRepository {
    if repoStore != nil {
        return repoStore.Users().WithContext(ctx)
    }
    if db != nil {
        return repository.NewUserRepository(db).WithContext(ctx)
    }
    return nil
}
