package repository

import (
	"context"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
)

type fileLogRepository struct {
	db *gorm.DB
}

func NewFileLogRepository(db *gorm.DB) domain.FileLogRepository {
	return &fileLogRepository{db: db}
}

func (r *fileLogRepository) Create(ctx context.Context, log *domain.FileLog) error {
	db := GetDB(ctx, r.db)
	return db.Create(log).Error
}

func (r *fileLogRepository) GetByTelegramID(ctx context.Context, telegramID int64, limit int, offset int) ([]domain.FileLog, error) {
	db := GetDB(ctx, r.db)
	var logs []domain.FileLog
	q := db.Where("telegram_id = ?", telegramID).Order("id DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Find(&logs).Error
	return logs, err
}

func (r *fileLogRepository) CountByTelegramID(ctx context.Context, telegramID int64) (int64, error) {
	db := GetDB(ctx, r.db)
	var count int64
	err := db.Model(&domain.FileLog{}).Where("telegram_id = ?", telegramID).Count(&count).Error
	return count, err
}

func (r *fileLogRepository) TotalDownloadedBytes(ctx context.Context, telegramID int64) (int64, error) {
	db := GetDB(ctx, r.db)
	var total *int64
	err := db.Model(&domain.FileLog{}).
		Where("telegram_id = ?", telegramID).
		Select("COALESCE(SUM(file_size_byte), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}
