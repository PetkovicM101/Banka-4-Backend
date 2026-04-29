package service

import (
	"context"
	"strings"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/repository"
)

type ProfitService struct {
	repo       repository.ProfitRepository
	userClient client.UserServiceClient
}

func NewProfitService(repo repository.ProfitRepository, userClient client.UserServiceClient) *ProfitService {
	return &ProfitService{
		repo:       repo,
		userClient: userClient,
	}
}
func (s *ProfitService) GetActuaryProfits(ctx context.Context) ([]dto.ActuaryProfitResponse, error) {

	actuaries, err := s.repo.GetAllActuaries(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.ActuaryProfitResponse

	for _, a := range actuaries {

		// user-service call per actuary
		emp, err := s.userClient.GetEmployeeById(ctx, uint64(a.EmployeeID))
		if err != nil {
			continue // ili log
		}

		nameParts := strings.Split(emp.FullName, " ")

		first := ""
		last := ""
		if len(nameParts) > 0 {
			first = nameParts[0]
		}
		if len(nameParts) > 1 {
			last = nameParts[1]
		}

		role := "agent"
		if emp.IsSupervisor {
			role = "supervisor"
		}

		result = append(result, dto.ActuaryProfitResponse{
			FirstName: first,
			LastName:  last,
			Role:      role,
			ProfitRSD: a.ProfitRSD, // ← iz trading baze, NE user service
		})
	}

	return result, nil
}
func (s *ProfitService) GetFundPositions(ctx context.Context) ([]dto.FundPositionResponse, error) {

	funds, err := s.repo.GetAllInvestmentFunds(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.FundPositionResponse

	for _, f := range funds {

		total := 0.0
		for _, p := range f.Positions {
			total += p.TotalInvestedAmount
		}

		bankPct := 10.0
		bankValue := total * bankPct / 100
		profit := bankValue * 0.2

		managerName := ""

		manager, err := s.userClient.GetEmployeeById(ctx, uint64(f.ManagerID))
		if err == nil {
			managerName = manager.FullName
		}

		result = append(result, dto.FundPositionResponse{
			FundName:       f.Name,
			ManagerName:    managerName,
			BankSharePct:   bankPct,
			BankShareValue: bankValue,
			Profit:         profit,
		})
	}

	return result, nil
}
