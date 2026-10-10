package service

import (
	"context"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type adminService struct {
	adminRepo domain.AdminRepository
}

func NewAdminService(adminRepo domain.AdminRepository) domain.AdminService {
	return &adminService{adminRepo: adminRepo}
}

func (s *adminService) IsAdmin(ctx context.Context, telegramID int64) (bool, error) {
	return s.adminRepo.IsAdminTelegramID(ctx, telegramID)
}

func (s *adminService) CreateAdmin(ctx context.Context, username, password string, telegramID *int64) (*domain.Admin, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &domain.Admin{
		Username:   username,
		Password:   string(hashedPassword),
		TelegramID: telegramID,
	}

	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
}
