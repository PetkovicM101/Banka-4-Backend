package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type ProfitRepository interface {
	GetAllInvestmentFunds(ctx context.Context) ([]model.InvestmentFund, error)
	GetAllActuaries(ctx context.Context) ([]model.Actuary, error)
}
