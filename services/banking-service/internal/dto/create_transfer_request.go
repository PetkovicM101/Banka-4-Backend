package dto

import (
	_ "time"
)

type CreateTransferRequest struct {
	PayerAccountNumber     string  `json:"payer_account_number"     binding:"required"`
	RecipientAccountNumber string  `json:"recipient_account_number" binding:"required"`
	Amount                 float64 `json:"start_amount"             binding:"required,min=0.01"`
}
