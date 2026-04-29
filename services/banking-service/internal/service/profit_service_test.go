package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/client"
)

type fakeProfitUserClient struct {
	err        error
	fullName   string
	supervisor bool
}

func (f *fakeProfitUserClient) GetClientByID(ctx context.Context, id uint) (*pb.GetClientByIdResponse, error) {
	if f.err != nil {
		return nil, f.err
	}

	return &pb.GetClientByIdResponse{
		FullName: f.fullName,
	}, nil
}

func (f *fakeProfitUserClient) GetEmployeeByID(ctx context.Context, id uint) (*pb.GetEmployeeByIdResponse, error) {
	if f.err != nil {
		return nil, f.err
	}

	return &pb.GetEmployeeByIdResponse{
		FullName:     f.fullName,
		IsSupervisor: f.supervisor,
	}, nil
}

// ─────────────────────────────────────────────
// TRADING CLIENT MOCK (OK)
// ─────────────────────────────────────────────

type fakeTradingClient struct {
	actuaries   []client.Actuary
	funds       []client.FundPosition
	actErr      error
	fundsErr    error
	profitValue float64
}

func (f *fakeTradingClient) GetActuaries(ctx context.Context, authHeader string) ([]client.Actuary, error) {
	return f.actuaries, f.actErr
}

func (f *fakeTradingClient) GetActuaryProfit(ctx context.Context, id uint, authHeader string) (float64, error) {
	return f.profitValue, nil
}

func (f *fakeTradingClient) GetInvestmentFunds(ctx context.Context, authHeader string) ([]client.FundPosition, error) {
	return f.funds, f.fundsErr
}

// ─────────────────────────────────────────────
// TEST SUCCESS
// ─────────────────────────────────────────────

func TestGetBankProfit_Success(t *testing.T) {
	userClient := &fakeProfitUserClient{
		fullName:   "John Doe",
		supervisor: true,
	}

	tradingClient := &fakeTradingClient{
		actuaries: []client.Actuary{
			{ID: 1},
		},
		funds: []client.FundPosition{
			{
				FundName:       "Fund A",
				ManagerName:    "Manager X",
				BankSharePct:   10,
				BankShareValue: 1000,
				Profit:         500,
			},
		},
		profitValue: 1234,
	}

	svc := &ProfitService{
		userClient:    userClient,
		tradingClient: tradingClient,
	}

	resp, err := svc.GetBankProfit(context.Background(), "Bearer test")

	require.NoError(t, err)
	require.Len(t, resp.Actuaries, 1)
	require.Len(t, resp.Funds, 1)

	a := resp.Actuaries[0]
	require.Equal(t, "John", a.FirstName)
	require.Equal(t, "Doe", a.LastName)
	require.Equal(t, "supervisor", a.Role) // ✅ FIX

	f := resp.Funds[0]
	require.Equal(t, "Fund A", f.FundName)
}

// ─────────────────────────────────────────────
// ERROR TESTS
// ─────────────────────────────────────────────

func TestGetBankProfit_ActuariesError(t *testing.T) {
	tradingClient := &fakeTradingClient{
		actErr: fmt.Errorf("trading down"),
	}

	svc := &ProfitService{
		userClient:    &fakeProfitUserClient{},
		tradingClient: tradingClient,
	}

	resp, err := svc.GetBankProfit(context.Background(), "auth")

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestGetBankProfit_FundsError(t *testing.T) {
	tradingClient := &fakeTradingClient{
		actuaries: []client.Actuary{{ID: 1}},
		fundsErr:  fmt.Errorf("funds error"),
	}

	svc := &ProfitService{
		userClient:    &fakeProfitUserClient{},
		tradingClient: tradingClient,
	}

	resp, err := svc.GetBankProfit(context.Background(), "auth")

	require.Error(t, err)
	require.Nil(t, resp)
}
