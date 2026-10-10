package service

import (
	"context"
	"time"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

type subscriptionService struct {
	subRepo   domain.SubscriptionRepository
	txManager domain.TxManager
}

func NewSubscriptionService(subRepo domain.SubscriptionRepository, txManager domain.TxManager) domain.SubscriptionService {
	return &subscriptionService{
		subRepo:   subRepo,
		txManager: txManager,
	}
}

func getNextTehranMidnight() time.Time {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		loc = time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc).UTC()
}

func (s *subscriptionService) GetOrCreateActiveSubscription(ctx context.Context, telegramID int64) (*domain.Subscription, error) {
	// 1. Check for active paid subscription
	paidSub, err := s.subRepo.GetActivePaidByTelegramID(ctx, telegramID)
	if err == nil && paidSub != nil {
		if paidSub.IsExpired() {
			_ = s.subRepo.DeactivateByID(ctx, paidSub.ID)
		} else {
			return paidSub, nil
		}
	}

	// 2. Check for active free subscription
	freeSub, err := s.subRepo.GetActiveFreeByTelegramID(ctx, telegramID)
	if err == nil && freeSub != nil {
		if freeSub.IsExpired() {
			_ = s.subRepo.DeactivateByID(ctx, freeSub.ID)
		} else {
			return freeSub, nil
		}
	}

	// 3. No active subscription found or previous was expired; provision 2GB daily free until Tehran midnight
	expiresAt := getNextTehranMidnight()
	newFreeSub := &domain.Subscription{
		TelegramID:             telegramID,
		PlanType:               domain.PlanTypeFree,
		TrafficLimitBytes:      domain.DefaultDailyFreeTrafficBytes,
		TrafficUsedBytes:       0,
		ExtraTrafficLimitBytes: 0,
		StartedAt:              time.Now().UTC(),
		ExpiresAt:              expiresAt,
		IsActive:               true,
	}

	if err := s.subRepo.Create(ctx, newFreeSub); err != nil {
		return nil, err
	}

	return newFreeSub, nil
}

func (s *subscriptionService) ActivatePaidSubscription(ctx context.Context, telegramID int64, plan *domain.Plan, txID int64) (*domain.Subscription, error) {
	var newSub *domain.Subscription

	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Deactivate any currently active subscriptions (both free and previous paid)
		if err := s.subRepo.DeactivateByTelegramID(txCtx, telegramID); err != nil {
			return err
		}

		expiresAt := time.Now().UTC().AddDate(0, 0, plan.DurationDays)
		newSub = &domain.Subscription{
			TelegramID:             telegramID,
			PlanType:               domain.PlanTypePaid,
			TrafficLimitBytes:      plan.TrafficBytes,
			TrafficUsedBytes:       0,
			ExtraTrafficLimitBytes: 0,
			TransactionID:          &txID,
			StartedAt:              time.Now().UTC(),
			ExpiresAt:              expiresAt,
			IsActive:               true,
		}

		return s.subRepo.Create(txCtx, newSub)
	})

	if err != nil {
		return nil, err
	}

	return newSub, nil
}

func (s *subscriptionService) AddExtraTraffic(ctx context.Context, telegramID int64, pkg *domain.ExtraTrafficPackage) (*domain.Subscription, error) {
	paidSub, err := s.subRepo.GetActivePaidByTelegramID(ctx, telegramID)
	if err != nil || paidSub == nil || paidSub.IsExpired() {
		return nil, domain.ErrOnlyPaidCanAddTraffic
	}

	if err := s.subRepo.AddExtraTraffic(ctx, paidSub.ID, pkg.TrafficBytes); err != nil {
		return nil, err
	}

	paidSub.ExtraTrafficLimitBytes += pkg.TrafficBytes
	return paidSub, nil
}

func (s *subscriptionService) ConsumeTraffic(ctx context.Context, telegramID int64, bytes int64) (*domain.Subscription, error) {
	sub, err := s.GetOrCreateActiveSubscription(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if !sub.HasSufficientTraffic(bytes) {
		return nil, domain.ErrInsufficientQuota
	}

	if err := s.subRepo.IncrementUsedTraffic(ctx, sub.ID, bytes); err != nil {
		return nil, err
	}

	sub.TrafficUsedBytes += bytes
	return sub, nil
}

func (s *subscriptionService) GetSubscriptionStats(ctx context.Context, telegramID int64) (*domain.Subscription, int64, error) {
	sub, err := s.GetOrCreateActiveSubscription(ctx, telegramID)
	if err != nil {
		return nil, 0, err
	}
	return sub, sub.RemainingTrafficBytes(), nil
}
