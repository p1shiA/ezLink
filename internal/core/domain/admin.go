package domain

import "context"

type Admin struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string `gorm:"column:username;size:50;unique;not null" json:"username"`
	TelegramID *int64 `gorm:"column:telegram_id" json:"telegram_id,omitempty"`
	Password   string `gorm:"column:password;size:255;not null" json:"-"`
}

func (Admin) TableName() string {
	return "admins"
}

type AdminRepository interface {
	IsAdminTelegramID(ctx context.Context, telegramID int64) (bool, error)
	GetByUsername(ctx context.Context, username string) (*Admin, error)
	Create(ctx context.Context, admin *Admin) error
}

type AdminService interface {
	IsAdmin(ctx context.Context, telegramID int64) (bool, error)
	CreateAdmin(ctx context.Context, username, password string, telegramID *int64) (*Admin, error)
}
