package domain

import (
	"context"
	"time"
)

type FileLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TelegramID   int64     `gorm:"column:telegram_id;not null;index" json:"telegram_id"`
	ChannelID    int64     `gorm:"column:channel_id;not null" json:"channel_id"`
	FileID       string    `gorm:"column:file_id;not null" json:"file_id"`
	FileName     string    `gorm:"column:file_name" json:"file_name"`
	FileType     string    `gorm:"column:file_type" json:"file_type"`
	FileSizeByte int64     `gorm:"column:file_size_byte;not null" json:"file_size_byte"`
	Link         string    `gorm:"column:link;not null" json:"link"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (FileLog) TableName() string {
	return "file_logs"
}

type FileLogRepository interface {
	Create(ctx context.Context, log *FileLog) error
	GetByTelegramID(ctx context.Context, telegramID int64, limit int, offset int) ([]FileLog, error)
	CountByTelegramID(ctx context.Context, telegramID int64) (int64, error)
	TotalDownloadedBytes(ctx context.Context, telegramID int64) (int64, error)
}

type FileLogService interface {
	// ProcessFileLink checks subscription quota, deducts traffic 1:1, and saves the file log
	ProcessFileLink(ctx context.Context, telegramID, channelID int64, fileID, fileName, fileType string, fileSizeByte int64, link string) (*FileLog, *Subscription, error)
	GetUserFileLogs(ctx context.Context, telegramID int64, limit, offset int) ([]FileLog, error)
	GetUserStats(ctx context.Context, telegramID int64) (totalFiles int64, totalBytes int64, err error)
}
