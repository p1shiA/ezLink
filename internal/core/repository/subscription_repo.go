package repository

import (
	"context"
	"errors"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) domain.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	db := GetDB(ctx, r.db)
	return db.Create(sub).Error
}

func (r *subscriptionRepository) GetActiveByTelegramID(ctx context.Context, telegramID int64) (*domain.Subscription, error) {
	db := GetDB(ctx, r.db)
	var sub domain.Subscription
	// Order by plan_type descending ('paid' before 'free') and id descending
	err := db.Where("telegram_id = ? AND is_active = true", telegramID).
		Order("plan_type DESC, id DESC").
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNoActiveSubscription
		}
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) GetActivePaidByTelegramID(ctx context.Context, telegramID int64) (*domain.Subscription, error) {
	db := GetDB(ctx, r.db)
	var sub domain.Subscription
	err := db.Where("telegram_id = ? AND plan_type = ? AND is_active = true", telegramID, domain.PlanTypePaid).
		Order("id DESC").
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNoActiveSubscription
		}
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) GetActiveFreeByTelegramID(ctx context.Context, telegramID int64) (*domain.Subscription, error) {
	db := GetDB(ctx, r.db)
	var sub domain.Subscription
	err := db.Where("telegram_id = ? AND plan_type = ? AND is_active = true", telegramID, domain.PlanTypeFree).
		Order("id DESC").
		First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNoActiveSubscription
		}
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) DeactivateByTelegramID(ctx context.Context, telegramID int64) error {
	db := GetDB(ctx, r.db)
	return db.Model(&domain.Subscription{}).
		Where("telegram_id = ? AND is_active = true", telegramID).
		Update("is_active", false).Error
}

func (r *subscriptionRepository) DeactivateByID(ctx context.Context, id int64) error {
	db := GetDB(ctx, r.db)
	return db.Model(&domain.Subscription{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *subscriptionRepository) AddExtraTraffic(ctx context.Context, id int64, extraBytes int64) error {
	db := GetDB(ctx, r.db)
	return db.Model(&domain.Subscription{}).
		Where("id = ?", id).
		Update("extra_traffic_limit_bytes", gorm.Expr("extra_traffic_limit_bytes + ?", extraBytes)).Error
}

func (r *subscriptionRepository) IncrementUsedTraffic(ctx context.Context, id int64, bytes int64) error {
	db := GetDB(ctx, r.db)
	return db.Model(&domain.Subscription{}).
		Where("id = ?", id).
		Update("traffic_used_bytes", gorm.Expr("traffic_used_bytes + ?", bytes)).Error
}
