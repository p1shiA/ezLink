package repository

import (
	"context"
	"errors"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) domain.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) IsAdminTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	db := GetDB(ctx, r.db)
	var count int64
	err := db.Model(&domain.Admin{}).Where("telegram_id = ?", telegramID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *adminRepository) GetByUsername(ctx context.Context, username string) (*domain.Admin, error) {
	db := GetDB(ctx, r.db)
	var admin domain.Admin
	err := db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAdminNotFound
		}
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) Create(ctx context.Context, admin *domain.Admin) error {
	db := GetDB(ctx, r.db)
	return db.Create(admin).Error
}
