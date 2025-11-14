package repository

import (
    "context"

    "gorm.io/gorm"
)

type Store struct {
    db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
    return &Store{db: db}
}

func (s *Store) DB() *gorm.DB {
    return s.db
}

func (s *Store) Users() UserRepository {
    return NewUserRepository(s.db)
}

func (s *Store) Dishes() DishRepository {
    return NewDishRepository(s.db)
}

func (s *Store) Contents() ContentRepository {
    return NewContentRepository(s.db)
}

type TxFn func(ctx context.Context, txStore *Store) error

type TransactionManager struct {
    db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
    return &TransactionManager{db: db}
}

func (tm *TransactionManager) WithTransaction(ctx context.Context, fn TxFn) error {
    return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        txStore := &Store{db: tx}
        return fn(ctx, txStore)
    })
}
