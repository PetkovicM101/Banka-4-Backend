package repository

import (
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type listingRepository struct {
	db *gorm.DB
}

func NewListingRepository(db *gorm.DB) ListingRepository {
	return &listingRepository{db: db}
}

func (r *listingRepository) FindAll() ([]model.Listing, error) {
	var listings []model.Listing
	if err := r.db.Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

func (r *listingRepository) Upsert(listing *model.Listing) error {
	return r.db.
		Where(model.Listing{Ticker: listing.Ticker}).
		Assign(*listing).
		FirstOrCreate(listing).Error
}

func (r *listingRepository) UpdatePriceAndAsk(listing *model.Listing, price, ask float64) error {
	return r.db.Model(listing).Updates(model.Listing{
		Price: price,
		Ask:   ask,
	}).Error
}
func (r *listingRepository) Count(count *int64) error {
	return r.db.Model(&model.Listing{}).Count(count).Error
}
