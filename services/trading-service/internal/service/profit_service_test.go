package service

import (
	"context"
	"testing"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

//
// MOCKS
//

type MockProfitRepo struct {
	mock.Mock
}

func (m *MockProfitRepo) GetAllActuaries(ctx context.Context) ([]model.Actuary, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Actuary), args.Error(1)
}

func (m *MockProfitRepo) GetAllInvestmentFunds(ctx context.Context) ([]model.InvestmentFund, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.InvestmentFund), args.Error(1)
}

//
// fake user client
//

type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GetClientById(ctx context.Context, id uint64) (*pb.GetClientByIdResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*pb.GetClientByIdResponse), args.Error(1)
}

func (m *MockUserClient) GetClientByIdentityId(ctx context.Context, identityId uint64) (*pb.GetClientByIdResponse, error) {
	args := m.Called(ctx, identityId)
	return args.Get(0).(*pb.GetClientByIdResponse), args.Error(1)
}

func (m *MockUserClient) GetEmployeeById(ctx context.Context, id uint64) (*pb.GetEmployeeByIdResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*pb.GetEmployeeByIdResponse), args.Error(1)
}

func (m *MockUserClient) GetEmployeeByIdentityId(ctx context.Context, identityId uint64) (*pb.GetEmployeeByIdResponse, error) {
	args := m.Called(ctx, identityId)
	return args.Get(0).(*pb.GetEmployeeByIdResponse), args.Error(1)
}

func (m *MockUserClient) GetAllClients(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllClientsResponse, error) {
	args := m.Called(ctx, page, pageSize, firstName, lastName)
	return args.Get(0).(*pb.GetAllClientsResponse), args.Error(1)
}

func (m *MockUserClient) GetAllActuaries(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllActuariesResponse, error) {
	args := m.Called(ctx, page, pageSize, firstName, lastName)
	return args.Get(0).(*pb.GetAllActuariesResponse), args.Error(1)
}

func (m *MockUserClient) GetIdentityByUserId(ctx context.Context, userID uint64, userType string) (*pb.GetIdentityByUserIdResponse, error) {
	args := m.Called(ctx, userID, userType)
	return args.Get(0).(*pb.GetIdentityByUserIdResponse), args.Error(1)
}

//
// TEST 1 - ACTUARY PROFITS
//

func TestGetActuaryProfits(t *testing.T) {
	ctx := context.Background()

	repo := new(MockProfitRepo)
	userClient := new(MockUserClient)

	service := NewProfitService(repo, userClient)

	repo.On("GetAllActuaries", ctx).Return([]model.Actuary{
		{
			EmployeeID: 1,
			ProfitRSD:  1000,
		},
	}, nil)

	userClient.On("GetEmployeeById", ctx, uint64(1)).Return(&pb.GetEmployeeByIdResponse{
		FullName:     "Marko Markovic",
		IsSupervisor: true,
	}, nil)

	result, err := service.GetActuaryProfits(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	assert.Equal(t, "Marko", result[0].FirstName)
	assert.Equal(t, "Markovic", result[0].LastName)
	assert.Equal(t, "supervisor", result[0].Role)
	assert.Equal(t, float64(1000), result[0].ProfitRSD)
}

//
// TEST 2 - FUND POSITIONS
//

func TestGetFundPositions(t *testing.T) {
	ctx := context.Background()

	repo := new(MockProfitRepo)
	userClient := new(MockUserClient)

	service := NewProfitService(repo, userClient)

	repo.On("GetAllInvestmentFunds", ctx).Return([]model.InvestmentFund{
		{
			Name:      "Tech Fund",
			ManagerID: 10,
			Positions: []model.ClientFundPosition{
				{TotalInvestedAmount: 1000},
				{TotalInvestedAmount: 2000},
			},
		},
	}, nil)

	userClient.On("GetEmployeeById", ctx, uint64(10)).Return(&pb.GetEmployeeByIdResponse{
		FullName: "Ana Anic",
	}, nil)

	result, err := service.GetFundPositions(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// total = 3000
	// bank 10% = 300
	// profit = 60

	assert.Equal(t, "Tech Fund", result[0].FundName)
	assert.Equal(t, "Ana Anic", result[0].ManagerName)
	assert.Equal(t, 300.0, result[0].BankShareValue)
	assert.Equal(t, 60.0, result[0].Profit)
}
