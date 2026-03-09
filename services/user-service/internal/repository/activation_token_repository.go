package repository

import (
	"context"
	"user-service/internal/model"

	"gorm.io/gorm"
)

type ActivationTokenRepository interface {
	Create(ctx context.Context, token *model.ActivationToken) error
	FindByToken(ctx context.Context, token string) (*model.ActivationToken, error)
	Delete(ctx context.Context, token *model.ActivationToken) error
}

type activationTokenRepo struct {
	db *gorm.DB
}

func NewActivationTokenRepository(db *gorm.DB) ActivationTokenRepository {
	return &activationTokenRepo{db: db}
}

func (r *activationTokenRepo) Create(ctx context.Context, token *model.ActivationToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *activationTokenRepo) FindByToken(ctx context.Context, token string) (*model.ActivationToken, error) {
	var t model.ActivationToken
	result := r.db.WithContext(ctx).Where("token = ?", token).First(&t)
	if result.Error != nil {
		return nil, result.Error
	}
	return &t, nil
}

func (r *activationTokenRepo) Delete(ctx context.Context, token *model.ActivationToken) error {
	return r.db.WithContext(ctx).Delete(token).Error
}
