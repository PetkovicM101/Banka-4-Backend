package grpc

import (
	"context"
	stderrors "errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	apperrors "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/repository"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/service"
)

type BankingService struct {
	pb.UnimplementedBankingServiceServer
	accountRepo          repository.AccountRepository
	transactionRepo      repository.TransactionRepository
	transactionProcessor *service.TransactionProcessor
	exchangeService      service.CurrencyConverter
}

func NewBankingService(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	transactionProcessor *service.TransactionProcessor,
	exchangeService service.CurrencyConverter,
) *BankingService {
	return &BankingService{
		accountRepo:          accountRepo,
		transactionRepo:      transactionRepo,
		transactionProcessor: transactionProcessor,
		exchangeService:      exchangeService,
	}
}

func (s *BankingService) GetAccountByNumber(ctx context.Context, req *pb.GetAccountByNumberRequest) (*pb.GetAccountByNumberResponse, error) {
	account, err := s.accountRepo.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch account: %v", err)
	}

	if account == nil {
		return nil, status.Errorf(codes.NotFound, "account %s not found", req.AccountNumber)
	}

	return &pb.GetAccountByNumberResponse{
		AccountNumber:    account.AccountNumber,
		ClientId:         uint64(account.ClientID),
		AccountType:      string(account.AccountType),
		CurrencyCode:     string(account.Currency.Code),
		AvailableBalance: account.AvailableBalance,
	}, nil
}

func (s *BankingService) ExecuteTradeSettlement(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	fields := req.GetFields()

	sourceAccountNumber, err := getStringField(fields, "source_account_number")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	destinationAccountNumber, err := getStringField(fields, "destination_account_number")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	amount, err := getNumberField(fields, "amount")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	amountIsSource, err := getBoolField(fields, "amount_is_source")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than zero")
	}

	sourceAccount, err := s.accountRepo.FindByAccountNumber(ctx, sourceAccountNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch source account: %v", err)
	}
	if sourceAccount == nil {
		return nil, status.Errorf(codes.NotFound, "source account %s not found", sourceAccountNumber)
	}

	destinationAccount, err := s.accountRepo.FindByAccountNumber(ctx, destinationAccountNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch destination account: %v", err)
	}
	if destinationAccount == nil {
		return nil, status.Errorf(codes.NotFound, "destination account %s not found", destinationAccountNumber)
	}

	sourceAmount := amount
	destinationAmount := amount
	if sourceAccount.Currency.Code != destinationAccount.Currency.Code {
		if amountIsSource {
			destinationAmount, err = s.exchangeService.Convert(ctx, amount, sourceAccount.Currency.Code, destinationAccount.Currency.Code)
		} else {
			sourceAmount, err = s.exchangeService.Convert(ctx, amount, destinationAccount.Currency.Code, sourceAccount.Currency.Code)
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert currencies: %v", err)
		}
	}

	transaction := &model.Transaction{
		PayerAccountNumber:     sourceAccount.AccountNumber,
		RecipientAccountNumber: destinationAccount.AccountNumber,
		StartAmount:            sourceAmount,
		StartCurrencyCode:      sourceAccount.Currency.Code,
		EndAmount:              destinationAmount,
		EndCurrencyCode:        destinationAccount.Currency.Code,
		Status:                 model.TransactionProcessing,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	if err := s.transactionProcessor.ProcessTradeSettlement(ctx, transaction.TransactionID); err != nil {
		return nil, mapTradeSettlementError(err)
	}

	result, err := structpb.NewStruct(map[string]any{
		"transaction_id":            float64(transaction.TransactionID),
		"source_amount":             sourceAmount,
		"source_currency_code":      string(sourceAccount.Currency.Code),
		"destination_amount":        destinationAmount,
		"destination_currency_code": string(destinationAccount.Currency.Code),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build response: %v", err)
	}

	return result, nil
}

func getStringField(fields map[string]*structpb.Value, name string) (string, error) {
	value, ok := fields[name]
	if !ok {
		return "", fmt.Errorf("%s is required", name)
	}

	stringValue := value.GetStringValue()
	if stringValue == "" {
		return "", fmt.Errorf("%s must be a non-empty string", name)
	}

	return stringValue, nil
}

func getNumberField(fields map[string]*structpb.Value, name string) (float64, error) {
	value, ok := fields[name]
	if !ok {
		return 0, fmt.Errorf("%s is required", name)
	}

	switch kind := value.Kind.(type) {
	case *structpb.Value_NumberValue:
		return kind.NumberValue, nil
	default:
		return 0, fmt.Errorf("%s must be a number", name)
	}
}

func getBoolField(fields map[string]*structpb.Value, name string) (bool, error) {
	value, ok := fields[name]
	if !ok {
		return false, fmt.Errorf("%s is required", name)
	}

	switch kind := value.Kind.(type) {
	case *structpb.Value_BoolValue:
		return kind.BoolValue, nil
	default:
		return false, fmt.Errorf("%s must be a boolean", name)
	}
}

func mapTradeSettlementError(err error) error {
	var appErr *apperrors.AppError
	if stderrors.As(err, &appErr) {
		switch appErr.Code {
		case 400:
			return status.Error(codes.FailedPrecondition, appErr.Message)
		case 404:
			return status.Error(codes.NotFound, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	return status.Error(codes.Internal, err.Error())
}
