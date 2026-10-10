package service

import (
	"context"
	"time"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

type fileLogService struct {
	fileLogRepo domain.FileLogRepository
	subService  domain.SubscriptionService
	txManager   domain.TxManager
}

func NewFileLogService(fileLogRepo domain.FileLogRepository, subService domain.SubscriptionService, txManager domain.TxManager) domain.FileLogService {
	return &fileLogService{
		fileLogRepo: fileLogRepo,
		subService:  subService,
		txManager:   txManager,
	}
}

func (s *fileLogService) ProcessFileLink(
	ctx context.Context,
	telegramID, channelID int64,
	fileID, fileName, fileType string,
	fileSizeByte int64,
	link string,
) (*domain.FileLog, *domain.Subscription, error) {
	var fileLog *domain.FileLog
	var activeSub *domain.Subscription

	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Validate quota and deduct traffic 1:1 atomically
		sub, err := s.subService.ConsumeTraffic(txCtx, telegramID, fileSizeByte)
		if err != nil {
			return err
		}
		activeSub = sub

		// 2. Persist the file log
		fileLog = &domain.FileLog{
			TelegramID:   telegramID,
			ChannelID:    channelID,
			FileID:       fileID,
			FileName:     fileName,
			FileType:     fileType,
			FileSizeByte: fileSizeByte,
			Link:         link,
			CreatedAt:    time.Now().UTC(),
		}

		return s.fileLogRepo.Create(txCtx, fileLog)
	})

	if err != nil {
		return nil, nil, err
	}

	return fileLog, activeSub, nil
}

func (s *fileLogService) GetUserFileLogs(ctx context.Context, telegramID int64, limit, offset int) ([]domain.FileLog, error) {
	return s.fileLogRepo.GetByTelegramID(ctx, telegramID, limit, offset)
}

func (s *fileLogService) GetUserStats(ctx context.Context, telegramID int64) (int64, int64, error) {
	count, err := s.fileLogRepo.CountByTelegramID(ctx, telegramID)
	if err != nil {
		return 0, 0, err
	}
	totalBytes, err := s.fileLogRepo.TotalDownloadedBytes(ctx, telegramID)
	if err != nil {
		return 0, 0, err
	}
	return count, totalBytes, nil
}
