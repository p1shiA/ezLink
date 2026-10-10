package repository

import (
	"context"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
)

type txContextKey struct{}

// ContextWithDB injects an active transaction into context
func ContextWithDB(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// GetDB retrieves an active transaction from context if available, otherwise returns the default DB
func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return defaultDB.WithContext(ctx)
}

// GormTxManager coordinates atomic transactions
type GormTxManager struct {
	db *gorm.DB
}

func NewGormTxManager(db *gorm.DB) domain.TxManager {
	return &GormTxManager{db: db}
}

func (m *GormTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := ContextWithDB(ctx, tx)
		return fn(txCtx)
	})
}
