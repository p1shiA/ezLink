package repository

import (
	"context"
	"errors"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
)

type planRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) domain.PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) GetAllPlans(ctx context.Context) ([]domain.Plan, error) {
	db := GetDB(ctx, r.db)
	var plans []domain.Plan
	err := db.Order("id ASC").Find(&plans).Error
	return plans, err
}

func (r *planRepository) GetPlanByID(ctx context.Context, id int) (*domain.Plan, error) {
	db := GetDB(ctx, r.db)
	var plan domain.Plan
	err := db.Where("id = ?", id).First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPlanNotFound
		}
		return nil, err
	}
	return &plan, nil
}

func (r *planRepository) GetAllExtraPackages(ctx context.Context) ([]domain.ExtraTrafficPackage, error) {
	db := GetDB(ctx, r.db)
	var pkgs []domain.ExtraTrafficPackage
	err := db.Order("id ASC").Find(&pkgs).Error
	return pkgs, err
}

func (r *planRepository) GetExtraPackageByID(ctx context.Context, id int) (*domain.ExtraTrafficPackage, error) {
	db := GetDB(ctx, r.db)
	var pkg domain.ExtraTrafficPackage
	err := db.Where("id = ?", id).First(&pkg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrExtraPackageNotFound
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *planRepository) CreatePlan(ctx context.Context, plan *domain.Plan) error {
	db := GetDB(ctx, r.db)
	return db.Create(plan).Error
}

func (r *planRepository) CreateExtraPackage(ctx context.Context, pkg *domain.ExtraTrafficPackage) error {
	db := GetDB(ctx, r.db)
	return db.Create(pkg).Error
}
