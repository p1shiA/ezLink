package database

import (
	"context"

	"github.com/p1shiA/ezLink/internal/core/domain"
	"gorm.io/gorm"
)

// SeedDefaults populates default plans and extra packages if the tables are empty
func SeedDefaults(ctx context.Context, db *gorm.DB) error {
	// Seed plans if none exist
	var planCount int64
	if err := db.WithContext(ctx).Model(&domain.Plan{}).Count(&planCount).Error; err != nil {
		return err
	}

	if planCount == 0 {
		defaultPlans := []domain.Plan{
			{
				Name:         "پلن ماهانه ۱۰ گیگابایت",
				PriceToman:   30000,
				TrafficBytes: 10 * 1024 * 1024 * 1024,
				DurationDays: 30,
			},
			{
				Name:         "پلن ماهانه ۳۰ گیگابایت",
				PriceToman:   70000,
				TrafficBytes: 30 * 1024 * 1024 * 1024,
				DurationDays: 30,
			},
			{
				Name:         "پلن ماهانه ۵۰ گیگابایت",
				PriceToman:   100000,
				TrafficBytes: 50 * 1024 * 1024 * 1024,
				DurationDays: 30,
			},
		}

		for _, p := range defaultPlans {
			if err := db.WithContext(ctx).Create(&p).Error; err != nil {
				return err
			}
		}
	}

	// Seed extra traffic packages if none exist
	var extraCount int64
	if err := db.WithContext(ctx).Model(&domain.ExtraTrafficPackage{}).Count(&extraCount).Error; err != nil {
		return err
	}

	if extraCount == 0 {
		defaultExtras := []domain.ExtraTrafficPackage{
			{
				PriceToman:   15000,
				TrafficBytes: 5 * 1024 * 1024 * 1024,
			},
			{
				PriceToman:   25000,
				TrafficBytes: 10 * 1024 * 1024 * 1024,
			},
		}

		for _, ep := range defaultExtras {
			if err := db.WithContext(ctx).Create(&ep).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
