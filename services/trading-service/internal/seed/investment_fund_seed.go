package seed

import (
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

func InvestmentFunds(db *gorm.DB) error {
	now := time.Now()

	funds := []model.InvestmentFund{
		{
			Name:                "Alpha Growth Fund",
			Description:         "Fond fokusiran na IT sektor sa agresivnom strategijom rasta.",
			MinimumContribution: 1000.00,
			ManagerID:           3,
			AccountNumber:       "444000000000000010",
			CreatedAt:           now,
		},
		{
			Name:                "Beta Stable Fund",
			Description:         "Konzervativni fond fokusiran na stabilne prihode i obveznice.",
			MinimumContribution: 5000.00,
			ManagerID:           7,
			AccountNumber:       "444000000000000011",
			CreatedAt:           now,
		},
	}

	for _, f := range funds {
		if err := db.FirstOrCreate(&f, model.InvestmentFund{Name: f.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}
