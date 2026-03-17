package dto

type CreatePaymentRequest struct {
	RecipientName    string  `json:"recipient_name" binding:"required"`
	RecipientAccount string  `json:"recipient_account" binding:"required"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
	ReferenceNumber  string  `json:"reference_number"`
	PaymentCode      string  `json:"payment_code"`
	Purpose          string  `json:"purpose"`
	PayerAccount     string  `json:"payer_account" binding:"required"`
	Currency         string  `json:"currency" binding:"required"`
}

type VerifyPaymentRequest struct {
	Code string `json:"code" binding:"required"`
}
