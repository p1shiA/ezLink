package domain

import (
	"context"
	"time"
)

type User struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TelegramID int64     `gorm:"column:telegram_id;uniqueIndex;not null" json:"telegram_id"`
	Username   string    `gorm:"column:username;size:255;index" json:"username"`
	FirstName  string    `gorm:"column:first_name;size:255" json:"first_name"`
	LastName   string    `gorm:"column:last_name;size:255" json:"last_name"`
	IsBanned   bool      `gorm:"column:is_banned;default:false;index" json:"is_banned"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

type UserRepository interface {
	Upsert(ctx context.Context, user *User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	UpdateBanStatus(ctx context.Context, telegramID int64, isBanned bool) error
	CountTotal(ctx context.Context) (int64, error)
}

type UserService interface {
	GetOrRegister(ctx context.Context, telegramID int64, username, firstName, lastName string) (*User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	SetBanStatus(ctx context.Context, telegramID int64, isBanned bool) error
	IsBanned(ctx context.Context, telegramID int64) (bool, error)
}
