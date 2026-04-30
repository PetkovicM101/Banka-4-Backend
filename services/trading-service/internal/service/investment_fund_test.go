package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/stretchr/testify/require"
)

// ── Fake Fund Repo ────────────────────────────────────────────────

type fakeFundRepo struct {
	findByIDResult      *model.InvestmentFund
	findByIDErr         error
	findByNameResult    *model.InvestmentFund
	findByNameErr       error
	createErr           error
	created             *model.InvestmentFund
	findAllResult       []model.InvestmentFund
	findAllTotal        int64
	findAllErr          error
	findByManagerResult []model.InvestmentFund
	findByManagerErr    error
}

func (f *fakeFundRepo) FindByName(ctx context.Context, name string) (*model.InvestmentFund, error) {
	return f.findByNameResult, f.findByNameErr
}

func (f *fakeFundRepo) FindByID(ctx context.Context, id uint) (*model.InvestmentFund, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeFundRepo) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.InvestmentFund, error) {
	return nil, nil
}

func (f *fakeFundRepo) GetAllInvestmentFunds(ctx context.Context) ([]model.InvestmentFund, error) {
	return f.findAllResult, f.findAllErr
}

func (f *fakeFundRepo) FindAll(ctx context.Context, name, sortBy, sortDir string, page, pageSize int) ([]model.InvestmentFund, int64, error) {
	return f.findAllResult, f.findAllTotal, f.findAllErr
}

func (f *fakeFundRepo) FindByManagerID(ctx context.Context, managerID uint) ([]model.InvestmentFund, error) {
	return f.findByManagerResult, f.findByManagerErr
}

func (f *fakeFundRepo) Create(ctx context.Context, fund *model.InvestmentFund) error {
	if f.createErr != nil {
		return f.createErr
	}
	fund.FundID = 1
	f.created = fund
	return nil
}

// ── Fake ClientFundPosition Repo ──────────────────────────────────

type fakePositionRepo struct {
	findResult *model.ClientFundPosition
	findErr    error
	upsertErr  error
}

func (f *fakePositionRepo) FindByClientAndFund(ctx context.Context, clientID uint, ownerType model.OwnerType, fundID uint) (*model.ClientFundPosition, error) {
	return f.findResult, f.findErr
}

func (f *fakePositionRepo) Upsert(ctx context.Context, position *model.ClientFundPosition) error {
	return f.upsertErr
}

// ── Fake ClientFundInvestment Repo ────────────────────────────────

type fakeInvestmentRepo struct {
	createErr error
}

func (f *fakeInvestmentRepo) Create(ctx context.Context, investment *model.ClientFundInvestment) error {
	return f.createErr
}

func (f *fakeInvestmentRepo) FindByClientAndFund(ctx context.Context, clientID uint, ownerType model.OwnerType, fundID uint) ([]model.ClientFundInvestment, error) {
	return nil, nil
}

// ── Fake Fund Banking Client ──────────────────────────────────────

type fakeFundBankingClient struct {
	createdAccountNumber string
	createFundAccountErr error
	getAccountResult     *pb.GetAccountByNumberResponse
	tradeSettlementErr   error
}

func (f *fakeFundBankingClient) GetAccountByNumber(_ context.Context, _ string) (*pb.GetAccountByNumberResponse, error) {
	if f.getAccountResult != nil {
		return f.getAccountResult, nil
	}
	return nil, nil
}
func (f *fakeFundBankingClient) HasActiveLoan(_ context.Context, _ uint64) (*pb.HasActiveLoanResponse, error) {
	return nil, nil
}
func (f *fakeFundBankingClient) CreatePaymentWithoutVerification(_ context.Context, _ *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	return nil, nil
}
func (f *fakeFundBankingClient) GetAccountsByClientID(_ context.Context, _ uint64) (*pb.GetAccountsByClientIDResponse, error) {
	return nil, nil
}
func (f *fakeFundBankingClient) ConvertCurrency(_ context.Context, amount float64, _, _ string) (float64, error) {
	return amount, nil
}
func (f *fakeFundBankingClient) ExecuteTradeSettlement(_ context.Context, _, _ string, _ pb.TradeSettlementDirection, _ float64) (*pb.ExecuteTradeSettlementResponse, error) {
	if f.tradeSettlementErr != nil {
		return nil, f.tradeSettlementErr
	}
	return &pb.ExecuteTradeSettlementResponse{}, nil
}
func (f *fakeFundBankingClient) GetAccountCurrency(_ context.Context, _ string) (string, error) {
	return "RSD", nil
}
func (f *fakeFundBankingClient) CreateFundAccount(_ context.Context, _ string, _ uint64) (string, error) {
	return f.createdAccountNumber, f.createFundAccountErr
}

type fakeFundUserClient struct {
	// configurable responses
	getClientByIdResp          *pb.GetClientByIdResponse
	getClientByIdErr           error

	getEmployeeByIdResp        *pb.GetEmployeeByIdResponse
	getEmployeeByIdErr         error

	getAllClientsResp          *pb.GetAllClientsResponse
	getAllClientsErr           error

	getAllActuariesResp        *pb.GetAllActuariesResponse
	getAllActuariesErr         error

	getIdentityByUserIdResp    *pb.GetIdentityByUserIdResponse
	getIdentityByUserIdErr     error
}

func (f *fakeFundUserClient) GetClientById(_ context.Context, _ uint64) (*pb.GetClientByIdResponse, error) {
	return f.getClientByIdResp, f.getClientByIdErr
}

func (f *fakeFundUserClient) GetClientByIdentityId(_ context.Context, _ uint64) (*pb.GetClientByIdResponse, error) {
	return f.getClientByIdResp, f.getClientByIdErr
}

func (f *fakeFundUserClient) GetEmployeeById(_ context.Context, _ uint64) (*pb.GetEmployeeByIdResponse, error) {
	return f.getEmployeeByIdResp, f.getEmployeeByIdErr
}

func (f *fakeFundUserClient) GetEmployeeByIdentityId(_ context.Context, _ uint64) (*pb.GetEmployeeByIdResponse, error) {
	return f.getEmployeeByIdResp, f.getEmployeeByIdErr
}

func (f *fakeFundUserClient) GetAllClients(_ context.Context, _, _ int32, _, _ string) (*pb.GetAllClientsResponse, error) {
	return f.getAllClientsResp, f.getAllClientsErr
}

func (f *fakeFundUserClient) GetAllActuaries(_ context.Context, _, _ int32, _, _ string) (*pb.GetAllActuariesResponse, error) {
	return f.getAllActuariesResp, f.getAllActuariesErr
}

func (f *fakeFundUserClient) GetIdentityByUserId(_ context.Context, _ uint64, _ string) (*pb.GetIdentityByUserIdResponse, error) {
	return f.getIdentityByUserIdResp, f.getIdentityByUserIdErr
}

// ── Helpers ───────────────────────────────────────────────────────

func fundSupervisorCtx() context.Context {
	employeeID := uint(25)
	return auth.SetAuthOnContext(context.Background(), &auth.AuthContext{
		IdentityID:   200,
		IdentityType: auth.IdentityEmployee,
		EmployeeID:   &employeeID,
	})
}

func fundClientCtx() context.Context {
	clientID := uint(99)
	return auth.SetAuthOnContext(context.Background(), &auth.AuthContext{
		IdentityID:   1,
		IdentityType: auth.IdentityClient,
		ClientID:     &clientID,
	})
}

func validFundRequest() dto.CreateFundRequest {
	return dto.CreateFundRequest{
		Name:                "Alpha Growth Fund",
		Description:         "Fund focused on the IT sector.",
		MinimumContribution: 1000.00,
	}
}

func newTestFundService(
	fundRepo *fakeFundRepo,
	ownershipRepo *fakeAssetOwnershipRepo,
	listingRepo *fakeListingRepo,
	bankingClient *fakeFundBankingClient,
	userClient *fakeFundUserClient,
) *InvestmentFundService {
	return NewInvestmentFundService(fundRepo, &fakePositionRepo{}, &fakeInvestmentRepo{}, ownershipRepo, listingRepo, bankingClient, userClient)
}

// ── CreateFund tests ──────────────────────────────────────────────

func TestCreateFund_Success(t *testing.T) {
	fundRepo := &fakeFundRepo{}
	bankingClient := &fakeFundBankingClient{createdAccountNumber: "444000112345678901"}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, bankingClient, &fakeFundUserClient{})

	resp, err := svc.CreateFund(fundSupervisorCtx(), validFundRequest())

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "Alpha Growth Fund", resp.Name)
	require.Equal(t, "444000112345678901", resp.AccountNumber)
	require.Equal(t, uint(25), resp.ManagerID)
	require.Equal(t, 1000.00, resp.MinimumContribution)
	require.WithinDuration(t, time.Now(), resp.CreatedAt, 5*time.Second)
}

func TestCreateFund_Unauthenticated(t *testing.T) {
	svc := newTestFundService(&fakeFundRepo{}, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.CreateFund(context.Background(), validFundRequest())

	require.Error(t, err)
	require.Contains(t, err.Error(), "not authenticated")
}

func TestCreateFund_NotEmployee(t *testing.T) {
	svc := newTestFundService(&fakeFundRepo{}, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.CreateFund(fundClientCtx(), validFundRequest())

	require.Error(t, err)
	require.Contains(t, err.Error(), "only employees")
}

func TestCreateFund_DuplicateName(t *testing.T) {
	fundRepo := &fakeFundRepo{
		findByNameResult: &model.InvestmentFund{Name: "Alpha Growth Fund"},
	}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.CreateFund(fundSupervisorCtx(), validFundRequest())

	require.Error(t, err)
	require.Contains(t, err.Error(), "already taken")
}

func TestCreateFund_FindByNameRepoError(t *testing.T) {
	fundRepo := &fakeFundRepo{
		findByNameErr: errors.New("db error"),
	}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.CreateFund(fundSupervisorCtx(), validFundRequest())

	require.Error(t, err)
}

func TestCreateFund_BankingClientError(t *testing.T) {
	fundRepo := &fakeFundRepo{}
	bankingClient := &fakeFundBankingClient{
		createFundAccountErr: fmt.Errorf("banking service unavailable"),
	}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, bankingClient, &fakeFundUserClient{})

	_, err := svc.CreateFund(fundSupervisorCtx(), validFundRequest())

	require.Error(t, err)
}

func TestCreateFund_RepoCreateError(t *testing.T) {
	fundRepo := &fakeFundRepo{
		createErr: errors.New("db error"),
	}
	bankingClient := &fakeFundBankingClient{createdAccountNumber: "444000112345678901"}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, bankingClient, &fakeFundUserClient{})

	_, err := svc.CreateFund(fundSupervisorCtx(), validFundRequest())

	require.Error(t, err)
}

// ── GetAllFunds tests ─────────────────────────────────────────────

func TestGetAllFunds_Success(t *testing.T) {
	fund := model.InvestmentFund{
		FundID:              1,
		Name:                "Alpha Growth Fund",
		Description:         "IT sector fund",
		MinimumContribution: 1000.0,
		ManagerID:           25,
		AccountNumber:       "444000000000000001",
		Positions: []model.ClientFundPosition{
			{TotalInvestedAmount: 300.0},
		},
	}
	ownership := model.AssetOwnership{AssetID: 10, Amount: 2.0}
	listing := model.Listing{AssetID: 10, Price: 100.0}

	fundRepo := &fakeFundRepo{findAllResult: []model.InvestmentFund{fund}, findAllTotal: 1}
	ownershipRepo := &fakeAssetOwnershipRepo{ownerships: []model.AssetOwnership{ownership}}
	listingRepo := &fakeListingRepo{byAssetIDs: []model.Listing{listing}}
	bankingClient := &fakeFundBankingClient{
		getAccountResult: &pb.GetAccountByNumberResponse{AvailableBalance: 500.0},
	}
	svc := newTestFundService(fundRepo, ownershipRepo, listingRepo, bankingClient, &fakeFundUserClient{})

	resp, err := svc.GetAllFunds(fundSupervisorCtx(), dto.ListFundsQuery{Page: 1, PageSize: 10})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(1), resp.Total)
	require.Len(t, resp.Data, 1)
	// securitiesValue = 2.0 * 100.0 = 200.0
	// fundValue = 500 (liquid) + 200 (securities) = 700
	require.Equal(t, 700.0, resp.Data[0].FundValue)
	// profit = 700 - 300 (invested) = 400
	require.Equal(t, 400.0, resp.Data[0].Profit)
	require.Equal(t, 500.0, resp.Data[0].LiquidAssets)
	require.Equal(t, 1, resp.Page)
	require.Equal(t, 10, resp.PageSize)
}

func TestGetAllFunds_Empty(t *testing.T) {
	fundRepo := &fakeFundRepo{findAllResult: []model.InvestmentFund{}, findAllTotal: 0}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	resp, err := svc.GetAllFunds(fundSupervisorCtx(), dto.ListFundsQuery{Page: 1, PageSize: 10})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(0), resp.Total)
	require.Empty(t, resp.Data)
}

func TestGetAllFunds_RepoError(t *testing.T) {
	fundRepo := &fakeFundRepo{findAllErr: errors.New("db error")}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.GetAllFunds(fundSupervisorCtx(), dto.ListFundsQuery{Page: 1, PageSize: 10})

	require.Error(t, err)
}

func TestGetAllFunds_OwnershipRepoError(t *testing.T) {
	fund := model.InvestmentFund{FundID: 1, Name: "Fund", AccountNumber: "444000000000000001"}
	fundRepo := &fakeFundRepo{findAllResult: []model.InvestmentFund{fund}, findAllTotal: 1}
	ownershipRepo := &fakeAssetOwnershipRepo{findErr: errors.New("db error")}
	svc := newTestFundService(fundRepo, ownershipRepo, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.GetAllFunds(fundSupervisorCtx(), dto.ListFundsQuery{Page: 1, PageSize: 10})

	require.Error(t, err)
}

// ── GetActuaryFunds tests ─────────────────────────────────────────

func TestGetActuaryFunds_Success(t *testing.T) {
	fund := model.InvestmentFund{
		FundID:        1,
		Name:          "Alpha Growth Fund",
		Description:   "IT sector fund",
		ManagerID:     25,
		AccountNumber: "444000000000000001",
	}
	ownership := model.AssetOwnership{AssetID: 5, Amount: 10.0}
	listing := model.Listing{AssetID: 5, Price: 50000.0}

	fundRepo := &fakeFundRepo{findByManagerResult: []model.InvestmentFund{fund}}
	ownershipRepo := &fakeAssetOwnershipRepo{ownerships: []model.AssetOwnership{ownership}}
	listingRepo := &fakeListingRepo{byAssetIDs: []model.Listing{listing}}
	bankingClient := &fakeFundBankingClient{
		getAccountResult: &pb.GetAccountByNumberResponse{AvailableBalance: 1500000.0},
	}
	svc := newTestFundService(fundRepo, ownershipRepo, listingRepo, bankingClient, &fakeFundUserClient{})

	resp, err := svc.GetActuaryFunds(fundSupervisorCtx(), 25)

	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, "Alpha Growth Fund", resp[0].Name)
	require.Equal(t, 1500000.0, resp[0].LiquidAssets)
	// securitiesValue = 10 * 50000 = 500000
	// fundValue = 1500000 + 500000 = 2000000
	require.Equal(t, 2000000.0, resp[0].FundValue)
}

func TestGetActuaryFunds_Empty(t *testing.T) {
	fundRepo := &fakeFundRepo{findByManagerResult: []model.InvestmentFund{}}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	resp, err := svc.GetActuaryFunds(fundSupervisorCtx(), 99)

	require.NoError(t, err)
	require.Empty(t, resp)
}

func TestGetActuaryFunds_RepoError(t *testing.T) {
	fundRepo := &fakeFundRepo{findByManagerErr: errors.New("db error")}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.GetActuaryFunds(fundSupervisorCtx(), 25)

	require.Error(t, err)
}

func TestGetActuaryFunds_OwnershipRepoError(t *testing.T) {
	fund := model.InvestmentFund{FundID: 1, Name: "Fund", ManagerID: 25}
	fundRepo := &fakeFundRepo{findByManagerResult: []model.InvestmentFund{fund}}
	ownershipRepo := &fakeAssetOwnershipRepo{findErr: errors.New("db error")}
	svc := newTestFundService(fundRepo, ownershipRepo, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.GetActuaryFunds(fundSupervisorCtx(), 25)

	require.Error(t, err)
}
