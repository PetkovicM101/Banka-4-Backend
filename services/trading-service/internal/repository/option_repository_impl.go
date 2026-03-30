package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &optionRepository{db: db}
}

// Upsert inserts or updates an Option matched by ListingID.
func (r *optionRepository) Upsert(ctx context.Context, option *model.Option) error {
	return r.db.WithContext(ctx).
		Where(model.Option{ListingID: option.ListingID}).
		Assign(*option).
		FirstOrCreate(option).Error
}

func (r *optionRepository) FindByListingIDs(ctx context.Context, ids []uint) ([]model.Option, error) {
	var options []model.Option
	err := r.db.WithContext(ctx).
		Where("listing_id IN ?", ids).
		Find(&options).Error
	return options, err
}
