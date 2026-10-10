package repository

import (
	"context"
	"errors"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Upsert(ctx context.Context, user *domain.User) error {
	db := GetDB(ctx, r.db)
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "telegram_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"username",
			"first_name",
			"last_name",
		}),
	}).Create(user).Error
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	db := GetDB(ctx, r.db)
	var user domain.User
	err := db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateBanStatus(ctx context.Context, telegramID int64, isBanned bool) error {
	db := GetDB(ctx, r.db)
	res := db.Model(&domain.User{}).Where("telegram_id = ?", telegramID).Update("is_banned", isBanned)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) CountTotal(ctx context.Context) (int64, error) {
	db := GetDB(ctx, r.db)
	var count int64
	err := db.Model(&domain.User{}).Count(&count).Error
	return count, err
}
