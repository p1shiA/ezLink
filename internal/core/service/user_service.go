package service

import (
	"context"
	"time"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetOrRegister(ctx context.Context, telegramID int64, username, firstName, lastName string) (*domain.User, error) {
	user := &domain.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.userRepo.Upsert(ctx, user); err != nil {
		return nil, err
	}

	return s.userRepo.GetByTelegramID(ctx, telegramID)
}

func (s *userService) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	return s.userRepo.GetByTelegramID(ctx, telegramID)
}

func (s *userService) SetBanStatus(ctx context.Context, telegramID int64, isBanned bool) error {
	return s.userRepo.UpdateBanStatus(ctx, telegramID, isBanned)
}

func (s *userService) IsBanned(ctx context.Context, telegramID int64) (bool, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return false, err
	}
	return user.IsBanned, nil
}
