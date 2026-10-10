package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

type transactionService struct {
	txRepo      domain.TransactionRepository
	planRepo    domain.PlanRepository
	subRepo     domain.SubscriptionRepository
	subService  domain.SubscriptionService
	gateway     domain.PaymentGateway
	txManager   domain.TxManager
	callbackURL string
}

func NewTransactionService(
	txRepo domain.TransactionRepository,
	planRepo domain.PlanRepository,
	subRepo domain.SubscriptionRepository,
	subService domain.SubscriptionService,
	gateway domain.PaymentGateway,
	txManager domain.TxManager,
	callbackURL string,
) domain.TransactionService {
	return &transactionService{
		txRepo:      txRepo,
		planRepo:    planRepo,
		subRepo:     subRepo,
		subService:  subService,
		gateway:     gateway,
		txManager:   txManager,
		callbackURL: callbackURL,
	}
}

func (s *transactionService) InitiatePlanPurchase(ctx context.Context, telegramID int64, planID int) (*domain.Transaction, string, error) {
	plan, err := s.planRepo.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, "", err
	}

	description := fmt.Sprintf("خرید %s - کاربر %d", plan.Name, telegramID)
	authority, paymentURL, err := s.gateway.RequestPayment(ctx, plan.PriceToman, description, s.callbackURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to request payment: %w", err)
	}

	tx := &domain.Transaction{
		TelegramID:      telegramID,
		TransactionType: domain.TxTypePlanPurchase,
		Status:          domain.TxStatusPending,
		PlanID:          &planID,
		PriceToman:      plan.PriceToman,
		Authority:       authority,
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, "", fmt.Errorf("failed to save transaction: %w", err)
	}

	return tx, paymentURL, nil
}

func (s *transactionService) InitiateExtraTrafficPurchase(ctx context.Context, telegramID int64, extraPackageID int) (*domain.Transaction, string, error) {
	// Rule: Extra traffic can only be added when an active paid subscription exists
	paidSub, err := s.subRepo.GetActivePaidByTelegramID(ctx, telegramID)
	if err != nil || paidSub == nil || paidSub.IsExpired() {
		return nil, "", domain.ErrOnlyPaidCanAddTraffic
	}

	pkg, err := s.planRepo.GetExtraPackageByID(ctx, extraPackageID)
	if err != nil {
		return nil, "", err
	}

	gb := pkg.TrafficBytes / (1024 * 1024 * 1024)
	description := fmt.Sprintf("خرید بسته ترافیک اضافی %d گیگابایت - کاربر %d", gb, telegramID)
	authority, paymentURL, err := s.gateway.RequestPayment(ctx, pkg.PriceToman, description, s.callbackURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to request payment: %w", err)
	}

	tx := &domain.Transaction{
		TelegramID:      telegramID,
		TransactionType: domain.TxTypeExtraTrafficPurchase,
		Status:          domain.TxStatusPending,
		ExtraTrafficID:  &extraPackageID,
		PriceToman:      pkg.PriceToman,
		Authority:       authority,
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, "", fmt.Errorf("failed to save transaction: %w", err)
	}

	return tx, paymentURL, nil
}

func (s *transactionService) VerifyAndCompletePayment(ctx context.Context, authority, status string) (*domain.Transaction, *domain.Subscription, error) {
	tx, err := s.txRepo.GetByAuthority(ctx, authority)
	if err != nil {
		return nil, nil, err
	}

	if tx.Status == domain.TxStatusCompleted {
		return tx, nil, domain.ErrTransactionAlreadyProcessed
	}

	// If customer cancelled or payment gateway returned failure status
	if !strings.EqualFold(status, "OK") {
		_ = s.txRepo.UpdateStatus(ctx, tx.ID, domain.TxStatusFailed, nil)
		return tx, nil, domain.ErrPaymentVerificationFailed
	}

	// Verify payment with payment gateway
	refID, err := s.gateway.VerifyPayment(ctx, authority, tx.PriceToman)
	if err != nil {
		_ = s.txRepo.UpdateStatus(ctx, tx.ID, domain.TxStatusFailed, nil)
		return tx, nil, fmt.Errorf("%w: %v", domain.ErrPaymentVerificationFailed, err)
	}

	var activeSub *domain.Subscription

	// Atomic database transaction to mark completed and apply subscription changes
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.txRepo.UpdateStatus(txCtx, tx.ID, domain.TxStatusCompleted, &refID); err != nil {
			return err
		}

		if tx.TransactionType == domain.TxTypePlanPurchase && tx.PlanID != nil {
			plan, err := s.planRepo.GetPlanByID(txCtx, *tx.PlanID)
			if err != nil {
				return err
			}
			activeSub, err = s.subService.ActivatePaidSubscription(txCtx, tx.TelegramID, plan, tx.ID)
			return err
		}

		if tx.TransactionType == domain.TxTypeExtraTrafficPurchase && tx.ExtraTrafficID != nil {
			pkg, err := s.planRepo.GetExtraPackageByID(txCtx, *tx.ExtraTrafficID)
			if err != nil {
				return err
			}
			activeSub, err = s.subService.AddExtraTraffic(txCtx, tx.TelegramID, pkg)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	tx.Status = domain.TxStatusCompleted
	tx.RefID = &refID

	return tx, activeSub, nil
}

func (s *transactionService) GetByAuthority(ctx context.Context, authority string) (*domain.Transaction, error) {
	return s.txRepo.GetByAuthority(ctx, authority)
}

func (s *transactionService) GetUserTransactions(ctx context.Context, telegramID int64, limit int) ([]domain.Transaction, error) {
	return s.txRepo.GetByTelegramID(ctx, telegramID, limit)
}
