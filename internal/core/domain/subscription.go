package domain

import (
	"context"
	"time"
)

const (
	PlanTypeFree = "free"
	PlanTypePaid = "paid"

	// DefaultDailyFreeTrafficBytes is 2GB in bytes
	DefaultDailyFreeTrafficBytes int64 = 2 * 1024 * 1024 * 1024
)

type Subscription struct {
	ID                     int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TelegramID             int64      `gorm:"column:telegram_id;not null;index" json:"telegram_id"`
	PlanType               string     `gorm:"column:plan_type;size:20;not null;index" json:"plan_type"`
	TrafficLimitBytes      int64      `gorm:"column:traffic_limit_bytes;not null" json:"traffic_limit_bytes"`
	TrafficUsedBytes       int64      `gorm:"column:traffic_used_bytes;default:0" json:"traffic_used_bytes"`
	ExtraTrafficLimitBytes int64      `gorm:"column:extra_traffic_limit_bytes;default:0" json:"extra_traffic_limit_bytes"`
	TransactionID          *int64     `gorm:"column:transaction_id" json:"transaction_id,omitempty"`
	StartedAt              time.Time  `gorm:"column:started_at;default:CURRENT_TIMESTAMP" json:"started_at"`
	ExpiresAt              time.Time  `gorm:"column:expires_at" json:"expires_at"`
	IsActive               bool       `gorm:"column:is_active;default:true;index" json:"is_active"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

// RemainingTrafficBytes computes remaining traffic based on plan type and extra packages
func (s *Subscription) RemainingTrafficBytes() int64 {
	if s.PlanType == PlanTypeFree {
		rem := s.TrafficLimitBytes - s.TrafficUsedBytes
		if rem < 0 {
			return 0
		}
		return rem
	}

	totalLimit := s.TrafficLimitBytes + s.ExtraTrafficLimitBytes
	rem := totalLimit - s.TrafficUsedBytes
	if rem < 0 {
		return 0
	}
	return rem
}

// HasSufficientTraffic checks if remaining bytes can cover required bytes
func (s *Subscription) HasSufficientTraffic(requiredBytes int64) bool {
	return s.RemainingTrafficBytes() >= requiredBytes
}

// IsExpired checks whether subscription expiry date has passed
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) error
	GetActiveByTelegramID(ctx context.Context, telegramID int64) (*Subscription, error)
	GetActivePaidByTelegramID(ctx context.Context, telegramID int64) (*Subscription, error)
	GetActiveFreeByTelegramID(ctx context.Context, telegramID int64) (*Subscription, error)
	DeactivateByTelegramID(ctx context.Context, telegramID int64) error
	DeactivateByID(ctx context.Context, id int64) error
	AddExtraTraffic(ctx context.Context, id int64, extraBytes int64) error
	IncrementUsedTraffic(ctx context.Context, id int64, bytes int64) error
}

type SubscriptionService interface {
	// GetOrCreateActiveSubscription returns current valid subscription (paid or free).
	// If paid is expired, deactivates it and falls back to daily free.
	// If free is missing or expired, lazy-renews it until midnight in Asia/Tehran.
	GetOrCreateActiveSubscription(ctx context.Context, telegramID int64) (*Subscription, error)
	ActivatePaidSubscription(ctx context.Context, telegramID int64, plan *Plan, txID int64) (*Subscription, error)
	AddExtraTraffic(ctx context.Context, telegramID int64, pkg *ExtraTrafficPackage) (*Subscription, error)
	ConsumeTraffic(ctx context.Context, telegramID int64, bytes int64) (*Subscription, error)
	GetSubscriptionStats(ctx context.Context, telegramID int64) (activeSub *Subscription, remainingBytes int64, err error)
}
