package service

import (
	"banking-service/internal/dto"
	"banking-service/internal/model"
	"banking-service/internal/repository"
)

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(req dto.CreatePaymentRequest) (*model.Payment, error) {

	// TODO: proveriti sredstva (#45)
	// TODO: proveriti limit
	// TODO: proveriti postojanje računa (#45)

	// TODO: currency conversion (#44)

	payment := &model.Payment{
		RecipientName:    req.RecipientName,
		RecipientAccount: req.RecipientAccount,
		Amount:           req.Amount,
		ReferenceNumber:  req.ReferenceNumber,
		PaymentCode:      req.PaymentCode,
		Purpose:          req.Purpose,
		PayerAccount:     req.PayerAccount,
		Currency:         req.Currency,
		Status:           model.PaymentProcessing,
	}

	err := s.repo.Create(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) VerifyPayment(id uint, code string) (*model.Payment, error) {

	payment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// TODO: mobile verification

	if code == "1234" {
		payment.Status = model.PaymentCompleted

		// TODO: save recipient
	} else {
		payment.Status = model.PaymentRejected
	}

	err = s.repo.Update(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
