package database

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager interface {
	WithinTransaction(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}

type gormTransactionManager struct {
	db *gorm.DB
}

func NewGORMTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTransactionManager{db: db}
}

type txKey struct{}

func (m *gormTransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func TxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}
