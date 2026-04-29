package client

import (
	"context"
)

type TradingClient interface {
	GetActuaries(ctx context.Context, authHeader string) ([]Actuary, error)
	GetActuaryProfit(ctx context.Context, actID uint, authHeader string) (float64, error)
	GetInvestmentFunds(ctx context.Context, authHeader string) ([]FundPosition, error)
}

type Actuary struct {
	ID uint
}

type FundPosition struct {
	FundName       string
	ManagerName    string
	BankSharePct   float64
	BankShareValue float64
	Profit         float64
}