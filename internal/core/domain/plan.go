package domain

import "context"

type Plan struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"column:name;size:50;unique;not null" json:"name"`
	PriceToman   int64  `gorm:"column:price_toman;not null" json:"price_toman"`
	TrafficBytes int64  `gorm:"column:traffic_bytes;not null" json:"traffic_bytes"`
	DurationDays int    `gorm:"column:duration_days;not null" json:"duration_days"`
}

func (Plan) TableName() string {
	return "plans"
}

type ExtraTrafficPackage struct {
	ID           int   `gorm:"primaryKey;autoIncrement" json:"id"`
	PriceToman   int64 `gorm:"column:price_toman;not null" json:"price_toman"`
	TrafficBytes int64 `gorm:"column:traffic_bytes;not null" json:"traffic_bytes"`
}

func (ExtraTrafficPackage) TableName() string {
	return "extra_traffic_packages"
}

type PlanRepository interface {
	GetAllPlans(ctx context.Context) ([]Plan, error)
	GetPlanByID(ctx context.Context, id int) (*Plan, error)
	GetAllExtraPackages(ctx context.Context) ([]ExtraTrafficPackage, error)
	GetExtraPackageByID(ctx context.Context, id int) (*ExtraTrafficPackage, error)
	CreatePlan(ctx context.Context, plan *Plan) error
	CreateExtraPackage(ctx context.Context, pkg *ExtraTrafficPackage) error
}

type PlanService interface {
	ListPlans(ctx context.Context) ([]Plan, error)
	GetPlanByID(ctx context.Context, id int) (*Plan, error)
	ListExtraPackages(ctx context.Context) ([]ExtraTrafficPackage, error)
	GetExtraPackageByID(ctx context.Context, id int) (*ExtraTrafficPackage, error)
}
