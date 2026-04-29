package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type InvestmentFundRepository interface {
	Create(ctx context.Context, fund *model.InvestmentFund) error
	FindByID(ctx context.Context, id uint) (*model.InvestmentFund, error)
	FindByAccountNumber(ctx context.Context, accountNumber string) (*model.InvestmentFund, error)
	FindByName(ctx context.Context, name string) (*model.InvestmentFund, error)
	FindHoldings(ctx context.Context, fundID uint) ([]model.AssetOwnership, error)
	GetAccountBalance(ctx context.Context, accountNumber string) (float64, error)
	CalculateTotalInvested(ctx context.Context, fundID uint) (float64, error)
	GetPerformanceHistory(ctx context.Context, fundID uint, limit int) ([]model.FundPerformance, error)
	SavePerformanceSnapshot(ctx context.Context, perf *model.FundPerformance) error
}
