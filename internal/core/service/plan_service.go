package service

import (
	"context"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

type planService struct {
	planRepo domain.PlanRepository
}

func NewPlanService(planRepo domain.PlanRepository) domain.PlanService {
	return &planService{planRepo: planRepo}
}

func (s *planService) ListPlans(ctx context.Context) ([]domain.Plan, error) {
	return s.planRepo.GetAllPlans(ctx)
}

func (s *planService) GetPlanByID(ctx context.Context, id int) (*domain.Plan, error) {
	return s.planRepo.GetPlanByID(ctx, id)
}

func (s *planService) ListExtraPackages(ctx context.Context) ([]domain.ExtraTrafficPackage, error) {
	return s.planRepo.GetAllExtraPackages(ctx)
}

func (s *planService) GetExtraPackageByID(ctx context.Context, id int) (*domain.ExtraTrafficPackage, error) {
	return s.planRepo.GetExtraPackageByID(ctx, id)
}
