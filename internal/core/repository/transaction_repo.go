package repository

import (
	"context"
	"errors"
	"time"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	db := GetDB(ctx, r.db)
	return db.Create(tx).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	db := GetDB(ctx, r.db)
	var tx domain.Transaction
	err := db.Preload("Plan").Preload("ExtraPackage").Where("id = ?", id).First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByAuthority(ctx context.Context, authority string) (*domain.Transaction, error) {
	db := GetDB(ctx, r.db)
	var tx domain.Transaction
	err := db.Preload("Plan").Preload("ExtraPackage").Where("authority = ?", authority).First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, id int64, status string, refID *string) error {
	db := GetDB(ctx, r.db)
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}
	if refID != nil {
		updates["ref_id"] = *refID
	}
	return db.Model(&domain.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *transactionRepository) GetByTelegramID(ctx context.Context, telegramID int64, limit int) ([]domain.Transaction, error) {
	db := GetDB(ctx, r.db)
	var txs []domain.Transaction
	q := db.Preload("Plan").Preload("ExtraPackage").Where("telegram_id = ?", telegramID).Order("id DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&txs).Error
	return txs, err
}
