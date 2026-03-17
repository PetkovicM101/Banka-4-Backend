package repository

import (
	"banking-service/internal/model"
	"context"
)

type TransferData struct {
	SourceAccountNum string
	DestAccountNum   string
	Amount           float64
	Description      string
}

type TransferHistory struct {
	TransactionID    *uint
	SourceAccountNum string
	DestAccountNum   string
	Amount           float64
	Description      string
	Status           string
	CreatedAt        string
}

type TransferRepository interface {
	// CreateTransfer zapisuje transfer u transaction tabelu (TransactionID će biti popunjen kasnije)
	CreateTransfer(ctx context.Context, transfer model.Transfer) error

	// GetTransferHistory vraća transfere za račun sa filteriranjem
	GetTransferHistory(ctx context.Context, accountNum string, status string, startDate, endDate string, page, pageSize int) ([]TransferHistory, int64, error)
}
