package service

import (
	"context"
	"strings"
	"sync"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/dto"
)

type ProfitService struct {
	userClient    client.UserClient
	tradingClient client.TradingClient
}

func NewProfitService(
	userClient client.UserClient,
	tradingClient client.TradingClient,
) *ProfitService {
	return &ProfitService{
		userClient:    userClient,
		tradingClient: tradingClient,
	}
}
func (s *ProfitService) GetBankProfit(ctx context.Context, authHeader string) (*dto.BankProfitResponse, error) {

	// 1. uzmi sve actuaries iz trading-service
	actuariesList, err := s.tradingClient.GetActuaries(ctx, authHeader)
	if err != nil {
		return nil, err
	}

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		actuaries []dto.ActuaryProfitResponse
	)

	// 2. paralelno fetch user + profit
	for _, a := range actuariesList {
		wg.Add(1)

		go func(actID uint) {
			defer wg.Done()

			emp, err := s.userClient.GetEmployeeByID(ctx, actID)
			if err != nil {
				return
			}

			profit, err := s.tradingClient.GetActuaryProfit(ctx, actID, authHeader)
			if err != nil {
				return
			}

			// split FullName -> First + Last
			firstName := ""
			lastName := ""

			if emp.FullName != "" {
				parts := strings.SplitN(emp.FullName, " ", 2)
				firstName = parts[0]
				if len(parts) > 1 {
					lastName = parts[1]
				}
			}

			// role mapping
			role := "agent"
			if emp.IsSupervisor {
				role = "supervisor"
			}

			mu.Lock()
			actuaries = append(actuaries, dto.ActuaryProfitResponse{
				FirstName: firstName,
				LastName:  lastName,
				Role:      role,
				ProfitRSD: profit,
			})
			mu.Unlock()

		}(a.ID)
	}

	wg.Wait()

	// 3. investment funds
	fundsData, err := s.tradingClient.GetInvestmentFunds(ctx, authHeader)
	if err != nil {
		return nil, err
	}

	funds := make([]dto.FundPositionResponse, 0, len(fundsData))

	for _, f := range fundsData {
		funds = append(funds, dto.FundPositionResponse{
			FundName:       f.FundName,
			ManagerName:    f.ManagerName,
			BankSharePct:   f.BankSharePct,
			BankShareValue: f.BankShareValue,
			ProfitRSD:      f.Profit,
		})
	}

	// 4. final response
	return &dto.BankProfitResponse{
		Actuaries: actuaries,
		Funds:     funds,
	}, nil
}
