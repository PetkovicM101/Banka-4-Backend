package repository

import (
	"banking-service/internal/model"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByID(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.First(&payment, id).Error
	return &payment, err
}

func (r *paymentRepository) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}
