package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type profitRepository struct {
	db *gorm.DB
}

func NewProfitRepository(db *gorm.DB) ProfitRepository {
	return &profitRepository{db: db}
}

func (r *profitRepository) GetAllActuaries(ctx context.Context) ([]model.Actuary, error) {
	var actuaries []model.Actuary

	err := r.db.WithContext(ctx).
		Find(&actuaries).Error

	return actuaries, err
}

func (r *profitRepository) GetAllInvestmentFunds(ctx context.Context) ([]model.InvestmentFund, error) {
	var funds []model.InvestmentFund

	err := r.db.WithContext(ctx).
		Preload("Positions").
		Find(&funds).Error

	return funds, err
}
