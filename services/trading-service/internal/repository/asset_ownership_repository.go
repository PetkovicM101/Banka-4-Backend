package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type AssetOwnershipRepository interface {
	FindByUserId(ctx context.Context, userId uint, ownerType model.OwnerType) ([]model.AssetOwnership, error)
	FindByID(ctx context.Context, id uint) (*model.AssetOwnership, error)
	Upsert(ctx context.Context, ownership *model.AssetOwnership) error
	FindAllPublic(ctx context.Context, page, pageSize int) ([]model.AssetOwnership, int64, error)
	UpdateOTCFields(ctx context.Context, ownershipID uint, publicAmount, reservedAmount float64) error
}
