package repository

import "banking-service/internal/model"

type PaymentRepository interface {
	Create(payment *model.Payment) error
	GetByID(id uint) (*model.Payment, error)
	Update(payment *model.Payment) error
}
