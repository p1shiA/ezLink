package domain

import (
	"context"
	"time"
)

const (
	TxTypePlanPurchase         = "plan_purchase"
	TxTypeExtraTrafficPurchase = "extra_traffic_purchase"

	TxStatusPending   = "pending"
	TxStatusCompleted = "completed"
	TxStatusFailed    = "failed"
)

type Transaction struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TelegramID      int64     `gorm:"column:telegram_id;not null;index" json:"telegram_id"`
	TransactionType string    `gorm:"column:transaction_type;size:50;not null" json:"transaction_type"`
	Status          string    `gorm:"column:status;size:20;default:pending;index" json:"status"`
	PlanID          *int      `gorm:"column:plan_id" json:"plan_id,omitempty"`
	ExtraTrafficID  *int      `gorm:"column:extra_traffic_id" json:"extra_traffic_id,omitempty"`
	PriceToman      int64     `gorm:"column:price_toman;not null" json:"price_toman"`
	Authority       string    `gorm:"column:authority;size:255;unique;not null" json:"authority"`
	RefID           *string   `gorm:"column:ref_id;size:255" json:"ref_id,omitempty"`
	CreatedAt       time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Associations for eager loading if needed
	Plan         *Plan                `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	ExtraPackage *ExtraTrafficPackage `gorm:"foreignKey:ExtraTrafficID" json:"extra_package,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	GetByAuthority(ctx context.Context, authority string) (*Transaction, error)
	UpdateStatus(ctx context.Context, id int64, status string, refID *string) error
	GetByTelegramID(ctx context.Context, telegramID int64, limit int) ([]Transaction, error)
}

type TransactionService interface {
	InitiatePlanPurchase(ctx context.Context, telegramID int64, planID int) (tx *Transaction, paymentURL string, err error)
	InitiateExtraTrafficPurchase(ctx context.Context, telegramID int64, extraPackageID int) (tx *Transaction, paymentURL string, err error)
	VerifyAndCompletePayment(ctx context.Context, authority, status string) (tx *Transaction, sub *Subscription, err error)
	GetByAuthority(ctx context.Context, authority string) (*Transaction, error)
	GetUserTransactions(ctx context.Context, telegramID int64, limit int) ([]Transaction, error)
}
