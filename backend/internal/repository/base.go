package repository

import (
	"context"
	"gorm.io/gorm"
)

type BaseRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db:  db,
		ctx: context.Background(),
	}
}

func (r *BaseRepository) WithContext(ctx context.Context) *BaseRepository {
	return &BaseRepository{
		db:  r.db,
		ctx: ctx,
	}
}

func (r *BaseRepository) WithTx(tx *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db:  tx,
		ctx: r.ctx,
	}
}

func (r *BaseRepository) DB() *gorm.DB {
	return r.db.WithContext(r.ctx)
}

type TransactionManagerImpl struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManagerImpl {
	return &TransactionManagerImpl{db: db}
}

func (tm *TransactionManagerImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}
