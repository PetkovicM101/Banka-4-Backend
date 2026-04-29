package repository

import (
	"context"
	"errors"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type investmentFundRepository struct {
	db *gorm.DB
}

func NewInvestmentFundRepository(db *gorm.DB) InvestmentFundRepository {
	return &investmentFundRepository{db: db}
}

func (r *investmentFundRepository) Create(ctx context.Context, fund *model.InvestmentFund) error {
	return r.db.WithContext(ctx).Create(fund).Error
}

func (r *investmentFundRepository) FindByID(ctx context.Context, id uint) (*model.InvestmentFund, error) {
	var fund model.InvestmentFund
	result := r.db.WithContext(ctx).First(&fund, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &fund, result.Error
}

func (r *investmentFundRepository) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.InvestmentFund, error) {
	var fund model.InvestmentFund
	result := r.db.WithContext(ctx).Where("account_number = ?", accountNumber).First(&fund)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &fund, result.Error
}

func (r *investmentFundRepository) FindByName(ctx context.Context, name string) (*model.InvestmentFund, error) {
	var fund model.InvestmentFund
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&fund)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &fund, result.Error
}

func (r *investmentFundRepository) FindHoldings(ctx context.Context, fundID uint) ([]model.AssetOwnership, error) {
	var holdings []model.AssetOwnership
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND owner_type = ?", fundID, model.OwnerTypeFund).
		Preload("Asset").
		Find(&holdings).Error
	return holdings, err
}

func (r *investmentFundRepository) GetAccountBalance(ctx context.Context, accountNumber string) (float64, error) {
	var balance float64
	err := r.db.WithContext(ctx).Table("accounts").
		Select("available_balance").
		Where("account_number = ?", accountNumber).
		Scan(&balance).Error
	return balance, err
}

func (r *investmentFundRepository) CalculateTotalInvested(ctx context.Context, fundID uint) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&model.AssetOwnership{}).
		Where("user_id = ? AND owner_type = ?", fundID, model.OwnerTypeFund).
		Select("COALESCE(SUM(amount * avg_buy_price_rsd), 0)").
		Scan(&total).Error
	return total, err
}

func (r *investmentFundRepository) GetPerformanceHistory(ctx context.Context, fundID uint, limit int) ([]model.FundPerformance, error) {
	var history []model.FundPerformance
	err := r.db.WithContext(ctx).
		Where("fund_id = ?", fundID).
		Order("date DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *investmentFundRepository) SavePerformanceSnapshot(ctx context.Context, perf *model.FundPerformance) error {
	return r.db.WithContext(ctx).Create(perf).Error
}
